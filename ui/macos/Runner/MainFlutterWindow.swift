import Cocoa
import FlutterMacOS

class MainFlutterWindow: NSWindow {
  override func awakeFromNib() {
    let flutterViewController = FlutterViewController()
    self.contentViewController = flutterViewController

    // 标题栏透明、内容铺到窗口顶部：侧边栏一直延伸到红绿灯下面，
    // 与「系统设置」「备忘录」等自带应用的窗口形态一致
    self.titlebarAppearsTransparent = true
    self.titleVisibility = .hidden
    self.styleMask.insert(.fullSizeContentView)

    self.minSize = NSSize(width: 760, height: 520)
    self.setContentSize(NSSize(width: 960, height: 640))
    self.center()
    // 记住用户调整过的窗口大小与位置
    self.setFrameAutosaveName("CopySyncMainWindow")

    RegisterGeneratedPlugins(registry: flutterViewController)

    super.awakeFromNib()
  }
}
