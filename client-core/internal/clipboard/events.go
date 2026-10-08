package clipboard

// 平台通过这里把系统事件交给后台服务。目前只有 Windows 会上报（见 windows.go），
// 其他平台这些值保持为空。

// changeNotify 在系统通知剪贴板变化时收到一个信号，Watcher 据此立即探测。
// 在 InitPlatform 里设置，之后不再改变。
var changeNotify chan struct{}

var (
	onSessionEnd func() // 用户注销或关机
	onResume     func() // 从睡眠中唤醒
)

// SetSystemHooks 注册系统事件的回调，须在 MainLoop 运行之前调用。
//
// sessionEnd 在系统即将结束会话时于主线程上调用，不能阻塞：返回之后进程随时可能被结束。
// resume 在唤醒后调用，此时网络多半已经变了。
func SetSystemHooks(sessionEnd, resume func()) {
	onSessionEnd, onResume = sessionEnd, resume
}
