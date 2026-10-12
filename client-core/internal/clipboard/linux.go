//go:build linux

package clipboard

// Linux 上经 X11 的 CLIPBOARD 选区读写剪贴板，纯 Go 实现（github.com/jezek/xgb），不依赖 xclip 之类的外部程序。
//
// Wayland 桌面（GNOME、KDE 等）都带着 XWayland，混成器会在 Wayland 剪贴板与 X11 选区之间双向同步，
// 所以同一套实现在 X11 和 Wayland 会话里都能用。连不上 X11 显示（没有 DISPLAY）时不可用。
//
// X11 的剪贴板不是一块共享内存，而是「谁复制谁持有」：
//   - 变化：用 XFixes 订阅 CLIPBOARD 的所有者变化，每变一次 changeCount 加一，并通知 Watcher 立即探测。
//   - 读取：请持有者把内容转换成某种格式（TARGETS、UTF8_STRING、image/png……）放到本机窗口的属性上；
//     大块数据按 ICCCM 的 INCR 协议分段传。
//   - 写入：本机成为持有者，在别的程序请求时按格式回应，同样支持 INCR。本机不再持有后（别人复制了新内容），
//     数据随之丢弃。所以本进程退出后，写进去的内容就粘贴不出来了，这与 X11 下所有程序的行为一致。

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"  // 注册解码器：别的程序可能只提供 GIF
	_ "image/jpeg" // 同上
	"image/png"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xfixes"
	"github.com/jezek/xgb/xproto"
)

// ErrNoDisplay 表示连不上 X11 显示：通常是后台服务不在图形会话里启动（比如经 SSH 启动）。
var ErrNoDisplay = errors.New("没有可用的 X11 显示（DISPLAY 未设置）")

const (
	convertTimeout = 2 * time.Second // 等持有者回应一次转换请求
	chunkTimeout   = 3 * time.Second // INCR 每一段之间最多等这么久
	reconnectEvery = 5 * time.Second
	// maxChunk 是单次写入属性的上限，超过就走 INCR。X 协议不开 BIG-REQUESTS 时单个请求最多 262140 字节，
	// 还要扣掉请求头；超长的请求会把整条连接写坏，所以留足余量
	maxChunk = 128 << 10
)

// 用到的格式名。text/plain 与 STRING 在不同程序里编码不一，读的时候都兜底按 UTF-8 解
const (
	tUTF8      = "UTF8_STRING"
	tPlainUTF8 = "text/plain;charset=utf-8"
	tPlain     = "text/plain"
	tText      = "TEXT"
	tString    = "STRING"
	tHTML      = "text/html"
	tPNG       = "image/png"
	tURIList   = "text/uri-list"
	tGnomeCopy = "x-special/gnome-copied-files" // GNOME 文件、Thunar 等文件管理器用它区分复制与剪切
	tKDECut    = "application/x-kde-cutselection"
	// tPasswordHint 是 KDE 定的约定，KeePassXC 等密码管理器复制密码时带着它：不要记录，也不要同步
	tPasswordHint = "x-kde-passwordManagerHint"
)

// imageTargets 是能读的图片格式，按优先顺序。
var imageTargets = []string{tPNG, "image/jpeg", "image/gif", "image/bmp"}

type linuxClipboard struct{}

func New() Clipboard { return linuxClipboard{} }

// InitPlatform 准备变化通知。连接 X11 推迟到第一次使用：图形会话可能还没就绪，之后会自动重试。
func InitPlatform() {
	changeNotify = make(chan struct{}, 1)
	x.mu.Lock()
	_ = x.connectLocked()
	x.mu.Unlock()
}

// X11 的事件在单独的 goroutine 里处理，主线程上不需要转动任何循环。
func pumpRunLoop(seconds float64) {
	time.Sleep(time.Duration(seconds * float64(time.Second)))
}

func (linuxClipboard) Permission() Permission        { return PermissionNotApplicable }
func (linuxClipboard) OpenPermissionSettings() error { return nil }

// ─────────────────────────── 连接 ───────────────────────────

// x 是本进程到 X 服务器的唯一连接。
var x = &x11{
	selNotify:  make(chan xproto.SelectionNotifyEvent, 4),
	propNotify: make(chan xproto.PropertyNotifyEvent, 64),
	sends:      make(map[sendKey]*incrSend),
}

type x11 struct {
	mu      sync.Mutex // 保护以下字段
	conn    *xgb.Conn
	win     xproto.Window
	atoms   map[string]xproto.Atom
	names   map[xproto.Atom]string
	own     *owned // 本机持有 CLIPBOARD 时提供的数据
	sends   map[sendKey]*incrSend
	lastTry time.Time
	lastErr error
	// cached 是最近一次 Peek 的结果，对应 cachedCount 那一次变化：没变化就不必再去问持有者
	cached      *Snapshot
	cachedCount int64

	// rw 让读取与取服务器时间串行进行：它们共用本机窗口上的属性和下面两个通道
	rw         sync.Mutex
	selNotify  chan xproto.SelectionNotifyEvent
	propNotify chan xproto.PropertyNotifyEvent

	count atomic.Int64 // 剪贴板变化次数，即 Snapshot.ChangeCount
}

// owned 是本机写入剪贴板的内容，按格式准备好。
type owned struct {
	time    xproto.Timestamp
	data    map[xproto.Atom][]byte
	order   []xproto.Atom // 回应 TARGETS 的顺序：越靠前越推荐
	snap    Snapshot
	content Content
}

type sendKey struct {
	win  xproto.Window
	prop xproto.Atom
}

// incrSend 是一次正在分段发送的 INCR 回应。
type incrSend struct {
	data    []byte
	typ     xproto.Atom
	off     int
	started time.Time
}

// connectLocked 在没有连接时连上 X 服务器。失败后 5 秒内不再重试。
func (x *x11) connectLocked() error {
	if x.conn != nil {
		return nil
	}
	if time.Since(x.lastTry) < reconnectEvery && x.lastErr != nil {
		return x.lastErr
	}
	x.lastTry = time.Now()
	x.lastErr = x.dialLocked()
	return x.lastErr
}

func (x *x11) dialLocked() error {
	display := os.Getenv("DISPLAY")
	if display == "" {
		return ErrNoDisplay
	}
	conn, err := xgb.NewConnDisplay(display)
	if err != nil {
		return fmt.Errorf("连接 X11 显示 %s: %w", display, err)
	}
	fail := func(what string, err error) error {
		conn.Close()
		return fmt.Errorf("%s: %w", what, err)
	}
	if err := xfixes.Init(conn); err != nil {
		return fail("X 服务器不支持 XFixes", err)
	}
	if _, err := xfixes.QueryVersion(conn, 5, 0).Reply(); err != nil {
		return fail("XFixes 版本协商", err)
	}

	screen := xproto.Setup(conn).DefaultScreen(conn)
	win, err := xproto.NewWindowId(conn)
	if err != nil {
		return fail("分配窗口", err)
	}
	// 不显示的输入窗口：只用来持有选区、接收转换结果。属性变化事件用于 INCR 和取服务器时间
	if err := xproto.CreateWindowChecked(conn, 0, win, screen.Root, 0, 0, 1, 1, 0,
		xproto.WindowClassInputOnly, 0, xproto.CwEventMask,
		[]uint32{xproto.EventMaskPropertyChange}).Check(); err != nil {
		return fail("创建窗口", err)
	}

	atoms := make(map[string]xproto.Atom)
	names := make(map[xproto.Atom]string)
	for _, name := range []string{
		"CLIPBOARD", "TARGETS", "TIMESTAMP", "INCR", "MULTIPLE", "SAVE_TARGETS", "CS_DATA", "CS_TIME",
		tUTF8, tPlainUTF8, tPlain, tText, tString, tHTML, tURIList, tGnomeCopy, tKDECut, tPasswordHint,
		tPNG, "image/jpeg", "image/gif", "image/bmp",
	} {
		r, err := xproto.InternAtom(conn, false, uint16(len(name)), name).Reply()
		if err != nil {
			return fail("登记格式名 "+name, err)
		}
		atoms[name], names[r.Atom] = r.Atom, name
	}

	if err := xfixes.SelectSelectionInputChecked(conn, win, atoms["CLIPBOARD"],
		xfixes.SelectionEventMaskSetSelectionOwner|xfixes.SelectionEventMaskSelectionWindowDestroy|
			xfixes.SelectionEventMaskSelectionClientClose).Check(); err != nil {
		return fail("订阅剪贴板变化", err)
	}

	x.conn, x.win, x.atoms, x.names = conn, win, atoms, names
	x.own, x.cached = nil, nil
	x.sends = make(map[sendKey]*incrSend)
	// 新连接：把它当作一次变化，Watcher 会重新探测一遍
	x.count.Add(1)
	go x.events(conn, win)
	return nil
}

// atom 返回格式名对应的原子，没登记过的就向服务器登记。
func (x *x11) atom(conn *xgb.Conn, name string) xproto.Atom {
	x.mu.Lock()
	a, ok := x.atoms[name]
	x.mu.Unlock()
	if ok {
		return a
	}
	r, err := xproto.InternAtom(conn, false, uint16(len(name)), name).Reply()
	if err != nil {
		return xproto.AtomNone
	}
	x.mu.Lock()
	x.atoms[name], x.names[r.Atom] = r.Atom, name
	x.mu.Unlock()
	return r.Atom
}

// name 返回原子对应的格式名，不认识的去服务器查。
func (x *x11) name(conn *xgb.Conn, a xproto.Atom) string {
	x.mu.Lock()
	n, ok := x.names[a]
	x.mu.Unlock()
	if ok {
		return n
	}
	r, err := xproto.GetAtomName(conn, a).Reply()
	if err != nil {
		return ""
	}
	x.mu.Lock()
	x.names[a] = r.Name
	x.mu.Unlock()
	return r.Name
}

// session 取出当前连接，必要时先连上。
func (x *x11) session() (*xgb.Conn, xproto.Window, error) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if err := x.connectLocked(); err != nil {
		return nil, 0, err
	}
	return x.conn, x.win, nil
}

// ─────────────────────────── 事件 ───────────────────────────

func (x *x11) events(conn *xgb.Conn, win xproto.Window) {
	for {
		ev, err := conn.WaitForEvent()
		if ev == nil && err == nil {
			x.lost(conn)
			return
		}
		if err != nil {
			// X 错误，比如请求方的窗口在回应前已经关了：不影响别的请求
			continue
		}
		switch e := ev.(type) {
		case xfixes.SelectionNotifyEvent:
			// 本机写入时已经计过数了（见 Write）
			if e.Owner != win {
				x.count.Add(1)
				notifyChange()
			}
		case xproto.SelectionRequestEvent:
			x.serve(conn, win, e)
		case xproto.SelectionClearEvent:
			x.mu.Lock()
			if x.conn == conn {
				x.own = nil
			}
			x.mu.Unlock()
		case xproto.SelectionNotifyEvent:
			select {
			case x.selNotify <- e:
			default:
			}
		case xproto.PropertyNotifyEvent:
			if e.Window == win {
				select {
				case x.propNotify <- e:
				default:
				}
			} else if e.State == xproto.PropertyDelete {
				x.continueSend(conn, e)
			}
		}
	}
}

// lost 处理连接断开。X 服务器没了几乎总是因为用户注销：按会话结束处理，正常退出；
// 下次登录时由开机自启重新拉起，而不是留下一个连不上新会话的后台服务。
func (x *x11) lost(conn *xgb.Conn) {
	x.mu.Lock()
	current := x.conn == conn
	if current {
		x.conn, x.own, x.cached = nil, nil, nil
		x.lastTry, x.lastErr = time.Now(), errors.New("与 X11 显示的连接已断开")
	}
	x.mu.Unlock()
	if current && onSessionEnd != nil {
		onSessionEnd()
	}
}

func notifyChange() {
	select {
	case changeNotify <- struct{}{}:
	default:
	}
}

// ─────────────────────────── 探测与读取 ───────────────────────────

func (linuxClipboard) Peek() (Snapshot, error) {
	conn, _, err := x.session()
	if err != nil {
		return Snapshot{}, err
	}
	count := x.count.Load()
	x.mu.Lock()
	if x.cached != nil && x.cachedCount == count {
		snap := *x.cached
		x.mu.Unlock()
		return snap, nil
	}
	own := x.own
	x.mu.Unlock()

	var snap Snapshot
	if own != nil {
		snap = own.snap
	} else {
		targets, err := x.targets(conn)
		if err != nil {
			return Snapshot{ChangeCount: count}, err
		}
		snap = classify(targets)
	}
	snap.ChangeCount = count

	x.mu.Lock()
	if x.count.Load() == count {
		x.cached, x.cachedCount = &snap, count
	}
	x.mu.Unlock()
	return snap, nil
}

// targets 问持有者能提供哪些格式。没有持有者时返回空列表。
func (x *x11) targets(conn *xgb.Conn) ([]string, error) {
	data, typ, err := x.convert(conn, "TARGETS")
	if errors.Is(err, errNoOwner) || errors.Is(err, errRefused) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if typ != xproto.AtomAtom && x.name(conn, typ) != "ATOM" {
		return nil, nil
	}
	var names []string
	for i := 0; i+4 <= len(data); i += 4 {
		if n := x.name(conn, xproto.Atom(xgb.Get32(data[i:]))); n != "" {
			names = append(names, n)
		}
	}
	return names, nil
}

// classify 按格式判断内容类型，规则与 Windows 一致：
// 文件优先；只有图片没有文本时才算图片（办公软件复制文字时会附带一张渲染出来的图）；再是 HTML、文本。
func classify(targets []string) Snapshot {
	has := make(map[string]bool, len(targets))
	for _, t := range targets {
		has[t] = true
	}
	hasText := has[tUTF8] || has[tPlainUTF8] || has[tPlain] || has[tString] || has[tText]
	hasImage := false
	for _, t := range imageTargets {
		hasImage = hasImage || has[t]
	}
	// 文件管理器复制文件时带着 x-special/gnome-copied-files 或 text/uri-list；浏览器复制网页内容时
	// 也可能带 text/uri-list，所以只有 uri-list、同时没有 HTML 和图片时才当作文件。真正是不是本地文件，读的时候再确认
	files := has[tGnomeCopy] || has[tKDECut] || (has[tURIList] && !has[tHTML] && !hasImage)

	snap := Snapshot{Types: targets, Sensitive: has[tPasswordHint]}
	switch {
	case files:
		snap.Kind, snap.ContentType = KindFile, "public.file-url"
	case hasImage && !hasText:
		snap.Kind, snap.ContentType = KindImage, "public.png"
	case has[tHTML]:
		snap.Kind, snap.ContentType = KindHTML, "public.html"
	case hasText:
		snap.Kind, snap.ContentType = KindText, "public.utf8-plain-text"
	case hasImage:
		snap.Kind, snap.ContentType = KindImage, "public.png"
	}
	return snap
}

func (linuxClipboard) Read() (Content, error) {
	conn, _, err := x.session()
	if err != nil {
		return Content{}, err
	}
	x.mu.Lock()
	own := x.own
	x.mu.Unlock()
	if own != nil {
		return own.content, nil
	}

	targets, err := x.targets(conn)
	if err != nil {
		return Content{}, err
	}
	has := make(map[string]bool, len(targets))
	for _, t := range targets {
		has[t] = true
	}
	kind := classify(targets).Kind

	if kind == KindFile {
		if files := x.readFiles(conn, has); len(files) > 0 {
			return Content{Kind: KindFile, Files: files}, nil
		}
		// 不是本地文件（比如网址）：按文本处理；没有文本格式时就用链接本身
		text := x.readText(conn, has)
		if text == "" && has[tURIList] {
			if data, _, err := x.convert(conn, tURIList); err == nil {
				text = strings.TrimSpace(strings.ReplaceAll(string(data), "\r\n", "\n"))
			}
		}
		return Content{Kind: KindText, Text: text}, nil
	}
	switch kind {
	case KindImage:
		for _, t := range imageTargets {
			if !has[t] {
				continue
			}
			data, _, err := x.convert(conn, t)
			if err != nil {
				return Content{}, err
			}
			img, err := toPNG(t, data)
			if err != nil {
				return Content{}, err
			}
			return Content{Kind: KindImage, Image: img}, nil
		}
	case KindHTML:
		data, _, err := x.convert(conn, tHTML)
		if err != nil {
			return Content{}, err
		}
		return Content{Kind: KindHTML, HTML: decodeHTML(data), Text: x.readText(conn, has)}, nil
	case KindText:
		return Content{Kind: KindText, Text: x.readText(conn, has)}, nil
	}
	return Content{}, nil
}

// readText 按 UTF8_STRING、text/plain;charset=utf-8、text/plain、STRING 的顺序取第一个有的。
func (x *x11) readText(conn *xgb.Conn, has map[string]bool) string {
	for _, t := range []string{tUTF8, tPlainUTF8, tPlain, tText, tString} {
		if !has[t] {
			continue
		}
		data, _, err := x.convert(conn, t)
		if err != nil {
			continue
		}
		return fromCRLF(decodeText(t, data))
	}
	return ""
}

// readFiles 取出本地文件的路径。只要还存在的本地文件，网址之类的都略过。
func (x *x11) readFiles(conn *xgb.Conn, has map[string]bool) []string {
	for _, t := range []string{tGnomeCopy, tURIList} {
		if !has[t] {
			continue
		}
		data, _, err := x.convert(conn, t)
		if err != nil {
			continue
		}
		if files := parseURIs(string(data)); len(files) > 0 {
			return files
		}
	}
	return nil
}

var (
	errNoOwner = errors.New("剪贴板为空")
	errRefused = errors.New("持有者拒绝了这种格式")
)

// convert 请持有者把剪贴板转换成 target 格式，返回数据和它声明的类型。
func (x *x11) convert(conn *xgb.Conn, target string) ([]byte, xproto.Atom, error) {
	x.rw.Lock()
	defer x.rw.Unlock()

	x.mu.Lock()
	win, clip, prop, incr := x.win, x.atoms["CLIPBOARD"], x.atoms["CS_DATA"], x.atoms["INCR"]
	x.mu.Unlock()
	tgt := x.atom(conn, target)

	owner, err := xproto.GetSelectionOwner(conn, clip).Reply()
	if err != nil {
		return nil, 0, err
	}
	if owner.Owner == xproto.AtomNone {
		return nil, 0, errNoOwner
	}

	drain(x.selNotify)
	drain(x.propNotify)
	xproto.DeleteProperty(conn, win, prop)
	xproto.ConvertSelection(conn, win, clip, tgt, prop, xproto.TimeCurrentTime)

	deadline := time.After(convertTimeout)
	for {
		select {
		case e := <-x.selNotify:
			if e.Target != tgt {
				continue // 之前某次超时的请求迟到的回应
			}
			if e.Property == xproto.AtomNone {
				return nil, 0, errRefused
			}
			data, typ, err := getProperty(conn, win, prop)
			if err != nil {
				return nil, 0, err
			}
			if typ == incr {
				data, typ, err = x.receiveIncr(conn, win, prop)
			}
			return data, typ, err
		case <-deadline:
			return nil, 0, fmt.Errorf("剪贴板持有者没有回应（格式 %s）", target)
		}
	}
}

// receiveIncr 按 INCR 协议分段接收：每删掉一次属性，持有者就写入下一段，写入空的一段表示结束。
func (x *x11) receiveIncr(conn *xgb.Conn, win xproto.Window, prop xproto.Atom) ([]byte, xproto.Atom, error) {
	var buf bytes.Buffer
	var typ xproto.Atom
	for {
		select {
		case e := <-x.propNotify:
			if e.Atom != prop || e.State != xproto.PropertyNewValue {
				continue
			}
			chunk, t, err := getProperty(conn, win, prop)
			if err != nil {
				return nil, 0, err
			}
			if t == xproto.AtomNone {
				// 属性已经不在了：这是持有者写入 INCR 标记时的那次事件，标记刚被读走删掉
				continue
			}
			if len(chunk) == 0 {
				return buf.Bytes(), typ, nil
			}
			typ = t
			buf.Write(chunk)
		case <-time.After(chunkTimeout):
			return nil, 0, errors.New("分段接收剪贴板内容超时")
		}
	}
}

// getProperty 读出并删掉本机窗口上的属性（删掉也是 INCR 协议里「要下一段」的信号）。
func getProperty(conn *xgb.Conn, win xproto.Window, prop xproto.Atom) ([]byte, xproto.Atom, error) {
	r, err := xproto.GetProperty(conn, true, win, prop, xproto.GetPropertyTypeAny, 0, 1<<28).Reply()
	if err != nil {
		return nil, 0, err
	}
	n := int(r.ValueLen) * int(r.Format) / 8
	if n > len(r.Value) {
		n = len(r.Value)
	}
	return r.Value[:n], r.Type, nil
}

func drain[T any](ch chan T) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}

// ─────────────────────────── 写入 ───────────────────────────

func (linuxClipboard) Write(content Content) error {
	conn, win, err := x.session()
	if err != nil {
		return err
	}
	own, err := x.prepare(conn, content)
	if err != nil {
		return err
	}
	t, err := x.serverTime(conn, win)
	if err != nil {
		return err
	}
	own.time = t

	x.mu.Lock()
	clip := x.atoms["CLIPBOARD"]
	x.own = own
	x.mu.Unlock()

	xproto.SetSelectionOwner(conn, win, clip, t)
	r, err := xproto.GetSelectionOwner(conn, clip).Reply()
	if err != nil || r.Owner != win {
		x.mu.Lock()
		x.own = nil
		x.mu.Unlock()
		return errors.New("没能成为剪贴板的持有者")
	}
	// Watcher 写入后马上 Peek，要看到新的计数；XFixes 随后报来的这次变化在 events 里略过
	x.count.Add(1)
	return nil
}

// prepare 把内容按各种格式准备好。
func (x *x11) prepare(conn *xgb.Conn, c Content) (*owned, error) {
	own := &owned{data: make(map[xproto.Atom][]byte), content: c}
	add := func(name string, data []byte) {
		a := x.atom(conn, name)
		if _, dup := own.data[a]; !dup {
			own.order = append(own.order, a)
		}
		own.data[a] = data
		own.snap.Types = append(own.snap.Types, name)
	}
	addText := func(s string) {
		b := []byte(s)
		add(tUTF8, b)
		add(tPlainUTF8, b)
		add(tPlain, b)
		add(tText, b)
		if isASCII(s) {
			add(tString, b) // 老程序按 Latin-1 解 STRING，只有纯 ASCII 才不会乱码
		}
	}

	switch c.Kind {
	case KindText:
		addText(c.Text)
		own.snap.Kind, own.snap.ContentType = KindText, "public.utf8-plain-text"
	case KindHTML:
		add(tHTML, []byte(c.HTML))
		if c.Text != "" {
			addText(c.Text)
		}
		own.snap.Kind, own.snap.ContentType = KindHTML, "public.html"
	case KindImage:
		if len(c.Image) == 0 {
			return nil, errors.New("没有要写入的图片数据")
		}
		add(tPNG, c.Image)
		own.snap.Kind, own.snap.ContentType = KindImage, "public.png"
	case KindFile:
		if len(c.Files) == 0 {
			return nil, errors.New("没有要写入的文件路径")
		}
		var uris []string
		for _, p := range c.Files {
			uris = append(uris, (&url.URL{Scheme: "file", Path: p}).String())
		}
		// 第一行 copy：文件管理器粘贴时复制而不是移动，缓存里的文件不能被挪走
		add(tGnomeCopy, []byte("copy\n"+strings.Join(uris, "\n")))
		add(tURIList, []byte(strings.Join(uris, "\r\n")+"\r\n"))
		add(tKDECut, []byte("0"))
		addText(strings.Join(c.Files, "\n"))
		own.snap.Kind, own.snap.ContentType = KindFile, "public.file-url"
	default:
		return nil, fmt.Errorf("不支持写入的类型：%v", c.Kind)
	}
	return own, nil
}

// serverTime 取 X 服务器的当前时间：成为持有者时要用真实的时间戳（ICCCM 不建议用 CurrentTime）。
// 做法是往本机窗口追加一个空属性，从随之而来的属性变化事件里读时间。
func (x *x11) serverTime(conn *xgb.Conn, win xproto.Window) (xproto.Timestamp, error) {
	x.rw.Lock()
	defer x.rw.Unlock()
	x.mu.Lock()
	prop := x.atoms["CS_TIME"]
	x.mu.Unlock()

	drain(x.propNotify)
	xproto.ChangeProperty(conn, xproto.PropModeAppend, win, prop, xproto.AtomString, 8, 0, nil)
	deadline := time.After(convertTimeout)
	for {
		select {
		case e := <-x.propNotify:
			if e.Atom == prop {
				return e.Time, nil
			}
		case <-deadline:
			return 0, errors.New("没有收到 X 服务器的时间")
		}
	}
}

// serve 回应别的程序对本机剪贴板内容的请求。
func (x *x11) serve(conn *xgb.Conn, win xproto.Window, e xproto.SelectionRequestEvent) {
	prop := e.Property
	if prop == xproto.AtomNone {
		prop = e.Target // 很老的程序不指定属性，约定用格式名本身
	}
	notify := xproto.SelectionNotifyEvent{
		Time: e.Time, Requestor: e.Requestor, Selection: e.Selection, Target: e.Target, Property: prop,
	}
	if !x.respond(conn, e, prop) {
		notify.Property = xproto.AtomNone
	}
	xproto.SendEvent(conn, false, e.Requestor, 0, string(notify.Bytes()))
}

func (x *x11) respond(conn *xgb.Conn, e xproto.SelectionRequestEvent, prop xproto.Atom) bool {
	x.mu.Lock()
	own, atoms := x.own, x.atoms
	x.mu.Unlock()
	if own == nil || e.Selection != atoms["CLIPBOARD"] {
		return false
	}

	switch e.Target {
	case atoms["TARGETS"]:
		list := append([]xproto.Atom{atoms["TARGETS"], atoms["TIMESTAMP"]}, own.order...)
		buf := make([]byte, 4*len(list))
		for i, a := range list {
			xgb.Put32(buf[4*i:], uint32(a))
		}
		xproto.ChangeProperty(conn, xproto.PropModeReplace, e.Requestor, prop, xproto.AtomAtom, 32,
			uint32(len(list)), buf)
		return true
	case atoms["TIMESTAMP"]:
		buf := make([]byte, 4)
		xgb.Put32(buf, uint32(own.time))
		xproto.ChangeProperty(conn, xproto.PropModeReplace, e.Requestor, prop, xproto.AtomInteger, 32, 1, buf)
		return true
	}

	data, ok := own.data[e.Target]
	if !ok {
		return false // MULTIPLE、SAVE_TARGETS 等都不支持，按约定回应「拒绝」
	}
	typ := e.Target
	if typ == atoms[tText] {
		typ = atoms[tUTF8] // TEXT 是「任意文本」，回应时要标明实际编码
	}
	if len(data) <= maxChunk {
		xproto.ChangeProperty(conn, xproto.PropModeReplace, e.Requestor, prop, typ, 8, uint32(len(data)), data)
		return true
	}

	// 太大：先放一个 INCR 属性写明总长度，请求方删掉它之后再一段段写入，见 continueSend
	x.mu.Lock()
	for k, s := range x.sends {
		if time.Since(s.started) > time.Minute {
			delete(x.sends, k) // 请求方中途不要了
		}
	}
	x.sends[sendKey{e.Requestor, prop}] = &incrSend{data: data, typ: typ, started: time.Now()}
	x.mu.Unlock()
	xproto.ChangeWindowAttributes(conn, e.Requestor, xproto.CwEventMask, []uint32{xproto.EventMaskPropertyChange})
	size := make([]byte, 4)
	xgb.Put32(size, uint32(len(data)))
	xproto.ChangeProperty(conn, xproto.PropModeReplace, e.Requestor, prop, atoms["INCR"], 32, 1, size)
	return true
}

// continueSend 在请求方删掉属性（读完上一段）后写入下一段。最后写入一段空数据表示结束。
func (x *x11) continueSend(conn *xgb.Conn, e xproto.PropertyNotifyEvent) {
	key := sendKey{e.Window, e.Atom}
	x.mu.Lock()
	s := x.sends[key]
	var chunk []byte
	if s != nil {
		end := min(s.off+maxChunk, len(s.data))
		chunk = s.data[s.off:end]
		s.off = end
		if len(chunk) == 0 {
			delete(x.sends, key)
		}
	}
	x.mu.Unlock()
	if s == nil {
		return
	}
	xproto.ChangeProperty(conn, xproto.PropModeReplace, e.Window, e.Atom, s.typ, 8, uint32(len(chunk)), chunk)
	if len(chunk) == 0 {
		xproto.ChangeWindowAttributes(conn, e.Window, xproto.CwEventMask, []uint32{0})
	}
}

// ─────────────────────────── 格式转换 ───────────────────────────

// decodeText 把文本格式的数据转成 UTF-8。STRING 按约定是 Latin-1，但不少程序其实放的是 UTF-8。
func decodeText(target string, data []byte) string {
	data = bytes.TrimRight(data, "\x00")
	if utf8.Valid(data) {
		return string(data)
	}
	runes := make([]rune, len(data))
	for i, b := range data {
		runes[i] = rune(b)
	}
	return string(runes)
}

// decodeHTML 解出 text/html。Firefox 等给的是带字节序标记的 UTF-16，其余是 UTF-8。
func decodeHTML(data []byte) string {
	switch {
	case len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE:
		return strings.TrimRight(decodeUTF16(data[2:], false), "\x00")
	case len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF:
		return strings.TrimRight(decodeUTF16(data[2:], true), "\x00")
	}
	return decodeText(tHTML, data)
}

func decodeUTF16(b []byte, bigEndian bool) string {
	u := make([]uint16, len(b)/2)
	for i := range u {
		if bigEndian {
			u[i] = uint16(b[2*i])<<8 | uint16(b[2*i+1])
		} else {
			u[i] = uint16(b[2*i]) | uint16(b[2*i+1])<<8
		}
	}
	return string(utf16.Decode(u))
}

// toPNG 把图片转成 PNG。BMP 去掉文件头后就是 Windows 的 DIB，复用那边的解码。
func toPNG(target string, data []byte) ([]byte, error) {
	if target == tPNG {
		return data, nil
	}
	var img image.Image
	var err error
	if target == "image/bmp" && len(data) > 14 && data[0] == 'B' && data[1] == 'M' {
		img, err = dibToImage(data[14:])
	} else {
		img, _, err = image.Decode(bytes.NewReader(data))
	}
	if err != nil {
		return nil, fmt.Errorf("解析 %s: %w", target, err)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// parseURIs 从 text/uri-list 或 x-special/gnome-copied-files 里取出还存在的本地文件路径。
func parseURIs(s string) []string {
	var out []string
	for _, line := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(strings.TrimRight(line, "\x00"))
		if line == "" || line == "copy" || line == "cut" || strings.HasPrefix(line, "#") {
			continue
		}
		u, err := url.Parse(line)
		if err != nil || u.Scheme != "file" || (u.Host != "" && u.Host != "localhost") || u.Path == "" {
			continue
		}
		if _, err := os.Stat(u.Path); err == nil {
			out = append(out, u.Path)
		}
	}
	return out
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}
