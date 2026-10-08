// Package sync 是同步引擎：把本机剪贴板变化广播出去，
// 并把对端送来的内容落盘、写入本机剪贴板。
//
// 传输模型（M-1 实测后确定）：按体积分流。
//   - 文本 / HTML：内容很小，直接内联在 offer 里，无需二次传输
//   - 图片、小于阈值的文件：随 offer 之后主动推送，对端可立即粘贴
//   - 超过阈值的文件：只发元数据，等用户在界面上点「拉取」
//
// 之所以不是"对方按下粘贴键才传"，是因为 macOS 的剪贴板延迟渲染
// 无法感知粘贴动作（详见 spikes/results.md）。
package sync

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"

	"github.com/baiyuze/copysync/client-core/internal/cache"
	"github.com/baiyuze/copysync/client-core/internal/clipboard"
	"github.com/baiyuze/copysync/client-core/internal/config"
	"github.com/baiyuze/copysync/client-core/internal/pack"
	"github.com/baiyuze/copysync/client-core/internal/store"
	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

const (
	// DataChannel 单条消息的上限。SCTP 对过大的消息支持参差，
	// 64 KiB 是各实现都稳妥的取值。
	chunkSize = 60 * 1024
	// 发送缓冲超过此值就暂停读取源文件，等它降下来再继续。
	// 没有背压的话，传大文件会把内存吃光。
	bufferHighWater = 1 << 20 // 1 MiB
	bufferLowWater  = 256 << 10

	// 预扫描的条目上限。超大目录全扫一遍会让复制动作明显卡顿，
	// 超限时只报概要，进度条退化为未知总量。
	maxScanEntries = 20000

	streamIdleTimeout = 2 * time.Minute
)

type Engine struct {
	log     *slog.Logger
	store   *store.Store
	cache   *cache.Cache
	watcher *clipboard.Watcher

	deviceID   string
	deviceName func() string
	loadConfig func() config.Config

	// 传输层
	send        func(deviceID string, msg *pb.PeerMessage) error
	broadcast   func(msg *pb.PeerMessage) int
	openStream  func(deviceID, label string) (*webrtc.DataChannel, error)
	onlinePeers func() []string

	// 事件回传给界面
	onRecord   func(*pb.ClipRecord)
	onProgress func(*pb.TransferProgress)

	mu sync.Mutex
	// 正在接收中的流：stream_id -> 接收器
	incoming map[string]*receiver

	// 每次传输上一次上报进度的时间，用于节流，见 emitProgress
	progressMu   sync.Mutex
	lastProgress map[string]time.Time
}

type Options struct {
	Logger      *slog.Logger
	Store       *store.Store
	Cache       *cache.Cache
	Watcher     *clipboard.Watcher
	DeviceID    string
	DeviceName  func() string
	LoadConfig  func() config.Config
	Send        func(deviceID string, msg *pb.PeerMessage) error
	Broadcast   func(msg *pb.PeerMessage) int
	OpenStream  func(deviceID, label string) (*webrtc.DataChannel, error)
	OnlinePeers func() []string
	OnRecord    func(*pb.ClipRecord)
	OnProgress  func(*pb.TransferProgress)
}

func New(opts Options) *Engine {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	return &Engine{
		log:          log,
		store:        opts.Store,
		cache:        opts.Cache,
		watcher:      opts.Watcher,
		deviceID:     opts.DeviceID,
		deviceName:   opts.DeviceName,
		loadConfig:   opts.LoadConfig,
		send:         opts.Send,
		broadcast:    opts.Broadcast,
		openStream:   opts.OpenStream,
		onlinePeers:  opts.OnlinePeers,
		onRecord:     opts.OnRecord,
		onProgress:   opts.OnProgress,
		incoming:     make(map[string]*receiver),
		lastProgress: make(map[string]time.Time),
	}
}

// ─────────────────────── 本机复制 → 广播 ───────────────────────

// HandleClipboardChange 在本机剪贴板变化时调用。
func (e *Engine) HandleClipboardChange(ctx context.Context, snap clipboard.Snapshot) {
	cfg := e.loadConfig()
	if !e.kindEnabled(cfg, snap.Kind) {
		e.log.Debug("该类型已在设置中关闭，跳过", "kind", snap.Kind)
		return
	}

	readCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	content, err := e.watcher.Read(readCtx)
	if err != nil {
		if errors.Is(err, clipboard.ErrPermissionRequired) {
			// 授权缺失是常态而非异常：界面会引导用户去开，这里安静跳过
			e.log.Debug("尚未获得剪贴板授权，跳过本次变化")
			return
		}
		e.log.Warn("读取剪贴板失败", "err", err)
		return
	}

	e.log.Debug("已读取剪贴板内容",
		"kind", content.Kind, "文本长度", len(content.Text),
		"HTML长度", len(content.HTML), "图片字节", len(content.Image), "文件", content.Files)

	clip, err := e.recordOutgoing(content, cfg)
	if err != nil {
		e.log.Warn("记录本机复制失败", "err", err)
		return
	}

	e.emitRecord(clip)

	// 没有在线对端时只留本地记录，不必尝试发送
	if len(e.onlinePeers()) == 0 {
		return
	}

	offer := e.buildOffer(clip, cfg)
	sent := e.broadcast(&pb.PeerMessage{
		Payload: &pb.PeerMessage_Offer{Offer: offer},
	})
	e.log.Info("已广播剪贴板内容",
		"clip", clip.ID, "kind", clip.Kind, "size", clip.TotalSize, "对端数", sent)

	// 小内容随即推送，对端无需等待即可粘贴
	if offer.GetWillPush() && clip.Kind == store.KindFile || offer.GetWillPush() && clip.Kind == store.KindImage {
		for _, peer := range e.onlinePeers() {
			go e.pushTo(context.WithoutCancel(ctx), peer, clip)
		}
	}
}

// recordOutgoing 把本机复制的内容落库。
func (e *Engine) recordOutgoing(c clipboard.Content, cfg config.Config) (store.Clip, error) {
	now := time.Now()
	clip := store.Clip{
		ID:               newID(),
		Status:           store.StatusReady,
		OriginDeviceID:   e.deviceID,
		OriginDeviceName: e.deviceName(),
		Outgoing:         true,
		CreatedAt:        now,
		ExpiresAt:        now.Add(cfg.HistoryTTL.Std()),
	}

	switch c.Kind {
	case clipboard.KindText, clipboard.KindHTML:
		clip.Kind = store.KindText
		clip.TextContent = c.Text
		if c.Kind == clipboard.KindHTML {
			clip.Kind = store.KindHTML
			clip.TextContent = c.HTML
			clip.TextPlain = c.Text
		}
		clip.TextPreview = preview(plainText(clip))
		clip.TotalSize = int64(len(clip.TextContent))

	case clipboard.KindImage:
		clip.Kind = store.KindImage
		clip.TotalSize = int64(len(c.Image))
		clip.TextPreview = fmt.Sprintf("图片 %s", humanBytes(clip.TotalSize))
		// 图片直接落盘，供对端拉取时复用
		rel, abs, err := e.cache.Dir(clip.ID)
		if err != nil {
			return clip, err
		}
		imgPath := filepath.Join(abs, "image.png")
		if err := os.WriteFile(imgPath, c.Image, 0o600); err != nil {
			return clip, err
		}
		clip.CachePath = rel
		clip.SourcePaths = []string{imgPath}
		clip.Items = []store.Item{{
			Name: "image.png", Size: clip.TotalSize, ContentType: "public.png",
		}}

	case clipboard.KindFile:
		clip.Kind = store.KindFile
		clip.SourcePaths = c.Files
		stats, err := pack.Scan(c.Files, maxScanEntries)
		if err != nil {
			e.log.Debug("预扫描未完成", "err", err)
		}
		clip.TotalSize = stats.Bytes
		for _, p := range c.Files {
			info, err := os.Stat(p)
			if err != nil {
				continue
			}
			clip.Items = append(clip.Items, store.Item{
				Name:  filepath.Base(p),
				Size:  info.Size(),
				IsDir: info.IsDir(),
			})
		}
		if len(clip.Items) == 0 {
			return clip, errors.New("剪贴板中的文件都不可访问")
		}
		clip.TextPreview = describeItems(clip.Items)

	default:
		return clip, fmt.Errorf("不支持的剪贴板类型: %v", c.Kind)
	}

	if err := e.store.PutClip(clip); err != nil {
		return clip, err
	}
	return clip, nil
}

// buildOffer 决定这条内容是随即推送还是等对方来取。
func (e *Engine) buildOffer(clip store.Clip, cfg config.Config) *pb.ClipOffer {
	offer := &pb.ClipOffer{
		ClipId:        clip.ID,
		Kind:          pb.ClipKind(clip.Kind),
		TotalSize:     clip.TotalSize,
		CreatedAtUnix: clip.CreatedAt.Unix(),
	}
	for _, it := range clip.Items {
		offer.Items = append(offer.Items, &pb.ClipItem{
			Name: it.Name, Size: it.Size, IsDir: it.IsDir, ContentType: it.ContentType,
		})
	}

	switch clip.Kind {
	case store.KindText, store.KindHTML:
		// 文本体积小，直接内联，省掉一次往返
		offer.TextContent = clip.TextContent
		offer.PlainText = clip.TextPlain
		offer.WillPush = false
	default:
		// 超过阈值的不主动推，避免无谓占用带宽和对端磁盘
		offer.WillPush = clip.TotalSize <= cfg.AutoSyncThresholdBytes
	}
	return offer
}

func (e *Engine) kindEnabled(cfg config.Config, k clipboard.Kind) bool {
	switch k {
	case clipboard.KindText:
		return cfg.SyncText
	case clipboard.KindHTML:
		return cfg.SyncHTML
	case clipboard.KindImage:
		return cfg.SyncImage
	case clipboard.KindFile:
		return cfg.SyncFile
	default:
		return false
	}
}

// ─────────────────────── 对端消息处理 ───────────────────────

func (e *Engine) HandlePeerMessage(ctx context.Context, from string, msg *pb.PeerMessage) {
	switch m := msg.GetPayload().(type) {
	case *pb.PeerMessage_Offer:
		e.handleOffer(ctx, from, m.Offer)
	case *pb.PeerMessage_Fetch:
		go e.handleFetch(context.WithoutCancel(ctx), from, m.Fetch.GetClipId())
	case *pb.PeerMessage_Reject:
		e.handleReject(m.Reject)
	case *pb.PeerMessage_Header:
		e.handleHeader(from, m.Header)
	case *pb.PeerMessage_Done:
		e.handleDone(m.Done)
	}
}

func (e *Engine) handleOffer(ctx context.Context, from string, offer *pb.ClipOffer) {
	cfg := e.loadConfig()
	now := time.Now()

	clip := store.Clip{
		ID:               offer.GetClipId(),
		Kind:             store.Kind(offer.GetKind()),
		OriginDeviceID:   from,
		OriginDeviceName: e.peerName(from),
		Outgoing:         false,
		TotalSize:        offer.GetTotalSize(),
		TextContent:      offer.GetTextContent(),
		TextPlain:        offer.GetPlainText(),
		CreatedAt:        time.Unix(offer.GetCreatedAtUnix(), 0),
		ExpiresAt:        now.Add(cfg.HistoryTTL.Std()),
	}
	for _, it := range offer.GetItems() {
		clip.Items = append(clip.Items, store.Item{
			Name: it.GetName(), Size: it.GetSize(),
			IsDir: it.GetIsDir(), ContentType: it.GetContentType(),
		})
	}

	switch {
	case clip.Kind == store.KindText || clip.Kind == store.KindHTML:
		// 文本随 offer 一起到了，立即可用
		clip.Status = store.StatusReady
		clip.TextPreview = preview(plainText(clip))
	case offer.GetWillPush():
		// 对端马上会推流过来，先记为接收中
		clip.Status = store.StatusFetching
		clip.TextPreview = describeItems(clip.Items)
	default:
		// 超过阈值，等用户在界面上点拉取
		clip.Status = store.StatusRemoteOnly
		clip.TextPreview = describeItems(clip.Items)
	}

	if err := e.store.PutClip(clip); err != nil {
		e.log.Warn("保存对端记录失败", "err", err)
		return
	}
	e.emitRecord(clip)
	e.log.Info("收到对端剪贴板内容",
		"clip", clip.ID, "kind", clip.Kind, "status", clip.Status, "from", from)

	// 文本类内容可以马上写进剪贴板
	if clip.Status == store.StatusReady && cfg.AutoApplyToClipboard {
		if err := e.applyToClipboard(ctx, clip); err != nil {
			e.log.Warn("写入剪贴板失败", "err", err)
		}
	}
}

// handleFetch 响应对端的拉取请求：校验源文件仍在，然后开流发送。
func (e *Engine) handleFetch(ctx context.Context, from, clipID string) {
	clip, err := e.store.GetClip(clipID)
	if err != nil {
		e.reject(from, clipID, "记录不存在")
		return
	}
	if err := e.pushTo(ctx, from, clip); err != nil {
		e.log.Warn("响应拉取失败", "clip", clipID, "err", err)
	}
}

func (e *Engine) handleReject(r *pb.FetchReject) {
	e.log.Warn("对端拒绝了拉取请求", "clip", r.GetClipId(), "原因", r.GetReason())
	_ = e.store.SetClipStatus(r.GetClipId(), store.StatusFailed, r.GetReason())
	if clip, err := e.store.GetClip(r.GetClipId()); err == nil {
		e.emitRecord(clip)
	}
}

func (e *Engine) reject(to, clipID, reason string) {
	_ = e.send(to, &pb.PeerMessage{
		Payload: &pb.PeerMessage_Reject{Reject: &pb.FetchReject{
			ClipId: clipID, Reason: reason,
		}},
	})
}

// ─────────────────────── 发送 ───────────────────────

// pushTo 把某条记录的内容流式发给指定对端。
func (e *Engine) pushTo(ctx context.Context, peer string, clip store.Clip) error {
	paths, err := e.sourcePaths(clip)
	if err != nil {
		e.reject(peer, clip.ID, err.Error())
		return err
	}

	streamID := newID()
	packaging := pb.Packaging_PACKAGING_TAR
	if len(paths) == 1 {
		if info, err := os.Stat(paths[0]); err == nil && !info.IsDir() {
			packaging = pb.Packaging_PACKAGING_RAW
		}
	}

	// 先发头，对端据此准备接收；随后的裸字节都属于这条流
	if err := e.send(peer, &pb.PeerMessage{
		Payload: &pb.PeerMessage_Header{Header: &pb.TransferHeader{
			ClipId:      clip.ID,
			StreamId:    streamID,
			TotalBytes:  clip.TotalSize,
			Compression: pb.Compression_COMPRESSION_ZSTD,
			Packaging:   packaging,
		}},
	}); err != nil {
		return fmt.Errorf("发送传输头: %w", err)
	}

	dc, err := e.openStream(peer, streamID)
	if err != nil {
		return fmt.Errorf("打开数据流: %w", err)
	}

	opened := make(chan struct{})
	dc.OnOpen(func() { close(opened) })
	select {
	case <-opened:
	case <-time.After(30 * time.Second):
		dc.Close()
		return errors.New("数据流打开超时")
	case <-ctx.Done():
		dc.Close()
		return ctx.Err()
	}

	start := time.Now()
	w := newChunkWriter(dc, func(sent int64) {
		e.emitProgress(clip.ID, sent, clip.TotalSize, start)
	})

	var packErr error
	if packaging == pb.Packaging_PACKAGING_RAW {
		packErr = packRaw(w, paths[0])
	} else {
		packErr = pack.Pack(w, paths)
	}
	if flushErr := w.Flush(); packErr == nil {
		packErr = flushErr
	}

	// 等发送缓冲排空再关闭，否则尾部数据会被丢弃
	w.WaitDrained(ctx, 60*time.Second)
	dc.Close()
	e.forgetProgress(clip.ID)

	ok := packErr == nil
	errMsg := ""
	if packErr != nil {
		errMsg = packErr.Error()
	}
	_ = e.send(peer, &pb.PeerMessage{
		Payload: &pb.PeerMessage_Done{Done: &pb.TransferDone{
			ClipId: clip.ID, StreamId: streamID, Ok: ok, Error: errMsg,
		}},
	})

	if packErr != nil {
		return packErr
	}
	e.log.Info("内容已发送",
		"clip", clip.ID, "to", peer, "字节", w.Sent(), "耗时", time.Since(start).Round(time.Millisecond))
	return nil
}

// sourcePaths 找出这条记录当前可用的源路径。
//
// 优先用原始位置（本机复制的内容），其次用缓存（从对端收来的内容）。
// 两者都没有时说明源文件已被移动或删除。
func (e *Engine) sourcePaths(clip store.Clip) ([]string, error) {
	if len(clip.SourcePaths) > 0 {
		var alive []string
		for _, p := range clip.SourcePaths {
			if _, err := os.Stat(p); err == nil {
				alive = append(alive, p)
			}
		}
		if len(alive) > 0 {
			return alive, nil
		}
	}
	if clip.CachePath != "" && e.cache.Exists(clip.CachePath) {
		return e.cache.Entries(clip.CachePath)
	}
	return nil, errors.New("源文件已不存在")
}

func packRaw(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	// 单文件也走 zstd，与 tar 分支保持一致的解码路径
	return pack.Pack(w, []string{path})
}

// ─────────────────────── 接收 ───────────────────────

// receiver 承接一条进行中的数据流。
type receiver struct {
	clipID   string
	streamID string
	from     string
	total    int64
	received int64
	pipeW    *io.PipeWriter
	done     chan error
	start    time.Time
}

func (e *Engine) handleHeader(from string, h *pb.TransferHeader) {
	rel, abs, err := e.cache.Dir(h.GetClipId())
	if err != nil {
		e.log.Warn("准备缓存目录失败", "err", err)
		return
	}

	pr, pw := io.Pipe()
	r := &receiver{
		clipID:   h.GetClipId(),
		streamID: h.GetStreamId(),
		from:     from,
		total:    h.GetTotalBytes(),
		pipeW:    pw,
		done:     make(chan error, 1),
		start:    time.Now(),
	}

	e.mu.Lock()
	e.incoming[h.GetStreamId()] = r
	e.mu.Unlock()

	_ = e.store.SetClipStatus(h.GetClipId(), store.StatusFetching, "")

	// 边收边解包：不等整条流收完就开始落盘，大文件才不会先堆在内存里
	go func() {
		tops, err := pack.Unpack(pr, abs)
		if err != nil {
			pr.CloseWithError(err)
			e.log.Warn("解包失败", "clip", r.clipID, "err", err)
			_ = e.store.SetClipStatus(r.clipID, store.StatusFailed, err.Error())
			r.done <- err
			return
		}
		e.finishReceive(r, rel, tops)
		r.done <- nil
	}()
}

// HandlePeerStream 接管对端新开的数据通道。
func (e *Engine) HandlePeerStream(from, label string, stream *p2p.Stream) {
	e.mu.Lock()
	r := e.incoming[label]
	e.mu.Unlock()

	if r == nil {
		// 头还没到就先收到了流。等一小会儿，控制通道与数据通道是两条独立的路，
		// 到达顺序并不保证。
		go func() {
			deadline := time.Now().Add(10 * time.Second)
			for time.Now().Before(deadline) {
				time.Sleep(50 * time.Millisecond)
				e.mu.Lock()
				r = e.incoming[label]
				e.mu.Unlock()
				if r != nil {
					e.bindStream(r, stream)
					return
				}
			}
			e.log.Warn("收到无对应传输头的数据流，已丢弃", "stream", label, "from", from)
			stream.Close()
		}()
		return
	}
	e.bindStream(r, stream)
}

func (e *Engine) bindStream(r *receiver, stream *p2p.Stream) {
	e.log.Debug("已绑定数据流", "stream", r.streamID, "clip", r.clipID,
		"state", stream.ReadyState())

	idle := time.AfterFunc(streamIdleTimeout, func() {
		r.pipeW.CloseWithError(errors.New("传输超时中断"))
		stream.Close()
	})

	// Stream 的队列保证顺序与完整性，这里只管按序搬进管道。
	// 队列在对端关闭数据流时关闭，循环随之退出并收尾。
	go func() {
		defer idle.Stop()
		defer r.pipeW.Close()

		for chunk := range stream.Chunks() {
			idle.Reset(streamIdleTimeout)
			if _, err := r.pipeW.Write(chunk); err != nil {
				stream.Close()
				return
			}
			r.received += int64(len(chunk))
			e.emitProgress(r.clipID, r.received, r.total, r.start)
		}
	}()
}

func (e *Engine) finishReceive(r *receiver, cacheRel string, tops []string) {
	e.mu.Lock()
	delete(e.incoming, r.streamID)
	e.mu.Unlock()
	e.forgetProgress(r.clipID)

	var size int64
	for _, p := range tops {
		if info, err := os.Stat(p); err == nil {
			size += info.Size()
		}
	}
	if err := e.store.SetClipCache(r.clipID, cacheRel, max64(size, r.received)); err != nil {
		e.log.Warn("更新缓存记录失败", "err", err)
	}

	clip, err := e.store.GetClip(r.clipID)
	if err != nil {
		return
	}
	e.emitRecord(clip)
	e.log.Info("内容接收完成",
		"clip", r.clipID, "from", r.from, "收到字节", r.received,
		"条目", len(tops), "耗时", time.Since(r.start).Round(time.Millisecond))

	if e.loadConfig().AutoApplyToClipboard {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := e.applyToClipboard(ctx, clip); err != nil {
			e.log.Warn("写入剪贴板失败", "err", err)
		}
	}
}

func (e *Engine) handleDone(d *pb.TransferDone) {
	if d.GetOk() {
		return
	}
	e.log.Warn("对端报告传输失败", "clip", d.GetClipId(), "err", d.GetError())
	_ = e.store.SetClipStatus(d.GetClipId(), store.StatusFailed, d.GetError())
	if clip, err := e.store.GetClip(d.GetClipId()); err == nil {
		e.emitRecord(clip)
	}
}

// ─────────────────────── 对外操作 ───────────────────────

// Fetch 主动拉取一条仅有元数据的记录。
func (e *Engine) Fetch(ctx context.Context, clipID string) error {
	clip, err := e.store.GetClip(clipID)
	if err != nil {
		return err
	}
	if clip.Outgoing {
		return errors.New("这是本机复制的内容，无需拉取")
	}
	if err := e.send(clip.OriginDeviceID, &pb.PeerMessage{
		Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: clipID}},
	}); err != nil {
		return fmt.Errorf("源设备当前不可达: %w", err)
	}
	_ = e.store.SetClipStatus(clipID, store.StatusFetching, "")
	if c, err := e.store.GetClip(clipID); err == nil {
		e.emitRecord(c)
	}
	return nil
}

// ApplyToClipboard 把某条已就绪的记录写入本机剪贴板。
func (e *Engine) ApplyToClipboard(ctx context.Context, clipID string) error {
	clip, err := e.store.GetClip(clipID)
	if err != nil {
		return err
	}
	return e.applyToClipboard(ctx, clip)
}

func (e *Engine) applyToClipboard(ctx context.Context, clip store.Clip) error {
	var content clipboard.Content

	switch clip.Kind {
	case store.KindText:
		content = clipboard.Content{Kind: clipboard.KindText, Text: clip.TextContent}
	case store.KindHTML:
		content = clipboard.Content{
			Kind: clipboard.KindHTML, HTML: clip.TextContent, Text: plainText(clip),
		}
	case store.KindImage, store.KindFile:
		paths, err := e.localPaths(clip)
		if err != nil {
			return err
		}
		if clip.Kind == store.KindImage && len(paths) == 1 {
			data, err := os.ReadFile(paths[0])
			if err == nil {
				content = clipboard.Content{Kind: clipboard.KindImage, Image: data}
				break
			}
		}
		content = clipboard.Content{Kind: clipboard.KindFile, Files: paths}
	default:
		return fmt.Errorf("不支持写入剪贴板的类型: %v", clip.Kind)
	}

	return e.watcher.Write(ctx, content)
}

// localPaths 返回这条记录在本机可用的实际路径。
func (e *Engine) localPaths(clip store.Clip) ([]string, error) {
	if clip.CachePath != "" && e.cache.Exists(clip.CachePath) {
		return e.cache.Entries(clip.CachePath)
	}
	if len(clip.SourcePaths) > 0 {
		var alive []string
		for _, p := range clip.SourcePaths {
			if _, err := os.Stat(p); err == nil {
				alive = append(alive, p)
			}
		}
		if len(alive) > 0 {
			return alive, nil
		}
	}
	return nil, errors.New("内容已不在本机（可能已过期被清理）")
}

// ─────────────────────── 辅助 ───────────────────────

func (e *Engine) emitRecord(clip store.Clip) {
	if e.onRecord == nil {
		return
	}
	items := make([]*pb.ClipItem, 0, len(clip.Items))
	for _, it := range clip.Items {
		items = append(items, &pb.ClipItem{
			Name: it.Name, Size: it.Size, IsDir: it.IsDir, ContentType: it.ContentType,
		})
	}
	e.onRecord(&pb.ClipRecord{
		Id:               clip.ID,
		Kind:             pb.ClipKind(clip.Kind),
		Status:           pb.ClipStatus(clip.Status),
		OriginDeviceId:   clip.OriginDeviceID,
		OriginDeviceName: clip.OriginDeviceName,
		Outgoing:         clip.Outgoing,
		Items:            items,
		TotalSize:        clip.TotalSize,
		TextPreview:      clip.TextPreview,
		CreatedAtUnix:    clip.CreatedAt.Unix(),
		ExpiresAtUnix:    clip.ExpiresAt.Unix(),
		Error:            clip.Error,
	})
}

// progressInterval 是同一次传输两次进度上报之间的最小间隔。
//
// 进度原本按数据块上报，一个大文件会产生上万条事件。事件队列每个订阅者只缓冲 64 条，
// 满了就丢弃新事件——「传输完成」那条记录更新也可能被挤掉，界面上的进度条就永远停不下来。
const progressInterval = 200 * time.Millisecond

func (e *Engine) emitProgress(clipID string, sent, total int64, start time.Time) {
	if e.onProgress == nil {
		return
	}
	done := total > 0 && sent >= total
	now := time.Now()
	e.progressMu.Lock()
	last, seen := e.lastProgress[clipID]
	if !done && seen && now.Sub(last) < progressInterval {
		e.progressMu.Unlock()
		return
	}
	e.lastProgress[clipID] = now
	e.progressMu.Unlock()
	elapsed := time.Since(start).Seconds()
	var rate float64
	if elapsed > 0 {
		rate = float64(sent) / elapsed
	}
	e.onProgress(&pb.TransferProgress{
		ClipId: clipID, Transferred: sent, Total: total, BytesPerSecond: rate,
	})
}

// forgetProgress 在一次传输结束后清掉它的节流记录。
func (e *Engine) forgetProgress(clipID string) {
	e.progressMu.Lock()
	delete(e.lastProgress, clipID)
	e.progressMu.Unlock()
}

// peerName 取对端设备名用于展示；查不到就退回 ID。
func (e *Engine) peerName(deviceID string) string {
	devices, err := e.store.ListDevices()
	if err != nil {
		return deviceID
	}
	for _, d := range devices {
		if d.ID == deviceID {
			return d.Name
		}
	}
	return deviceID
}

func newID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("clip-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func preview(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
	const limit = 120
	if len(s) <= limit {
		return s
	}
	// 按 rune 截断，避免把多字节字符切一半
	runes := []rune(s)
	if len(runes) <= limit {
		return s
	}
	return string(runes[:limit]) + "…"
}

func describeItems(items []store.Item) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0].Name
	default:
		return fmt.Sprintf("%s 等 %d 项", items[0].Name, len(items))
	}
}

func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for v := n / unit; v >= unit; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
