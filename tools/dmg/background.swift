// 绘制 DMG 安装窗口的背景图（1x 与 2x 两份，由 tiffutil 合成一张多分辨率 TIFF）。
//
//   swift tools/dmg/background.swift <输出目录>
//
// 版式与 scripts/package.sh 里 Finder 窗口的设置一一对应：
// 窗口 640×400，App 图标中心在 (160, 190)，「应用程序」在 (480, 190)。

import AppKit

let W: CGFloat = 640, H: CGFloat = 400
let ink = NSColor(srgbRed: 0x1D / 255.0, green: 0x1D / 255.0, blue: 0x1F / 255.0, alpha: 1)
let dim = NSColor(srgbRed: 0x6E / 255.0, green: 0x6E / 255.0, blue: 0x73 / 255.0, alpha: 1)
let faint = NSColor(srgbRed: 0xB4 / 255.0, green: 0xB4 / 255.0, blue: 0xB8 / 255.0, alpha: 1)
let paper = NSColor(srgbRed: 0xF5 / 255.0, green: 0xF5 / 255.0, blue: 0xF4 / 255.0, alpha: 1)

/// 以左上角为原点给出坐标，换算成 AppKit 的左下原点。
func y(_ top: CGFloat) -> CGFloat { H - top }

func text(_ s: String, size: CGFloat, weight: NSFont.Weight, color: NSColor, centerTop: CGFloat) {
    let attrs: [NSAttributedString.Key: Any] = [
        .font: NSFont.systemFont(ofSize: size, weight: weight),
        .foregroundColor: color,
    ]
    let str = NSAttributedString(string: s, attributes: attrs)
    let sz = str.size()
    str.draw(at: NSPoint(x: (W - sz.width) / 2, y: y(centerTop) - sz.height / 2))
}

func draw() {
    paper.setFill()
    NSRect(x: 0, y: 0, width: W, height: H).fill()

    // 两个图标之间的箭头：细线，不抢图标的戏
    let arrow = NSBezierPath()
    arrow.lineWidth = 2
    arrow.lineCapStyle = .round
    arrow.lineJoinStyle = .round
    arrow.move(to: NSPoint(x: 246, y: y(190)))
    arrow.line(to: NSPoint(x: 394, y: y(190)))
    arrow.move(to: NSPoint(x: 382, y: y(180)))
    arrow.line(to: NSPoint(x: 394, y: y(190)))
    arrow.line(to: NSPoint(x: 382, y: y(200)))
    faint.setStroke()
    arrow.stroke()

    text("把 CopySync 拖到「应用程序」文件夹", size: 14, weight: .medium, color: ink, centerTop: 300)
    text("Drag CopySync to the Applications folder", size: 12, weight: .regular, color: dim, centerTop: 322)
    text("首次打开如被拦截：系统设置 → 隐私与安全性 → 仍要打开", size: 11, weight: .regular,
         color: dim, centerTop: 366)
}

func render(scale: CGFloat, to url: URL) throws {
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: Int(W * scale),
                               pixelsHigh: Int(H * scale), bitsPerSample: 8, samplesPerPixel: 4,
                               hasAlpha: true, isPlanar: false, colorSpaceName: .deviceRGB,
                               bytesPerRow: 0, bitsPerPixel: 0)!
    rep.size = NSSize(width: W, height: H) // 2x 图的点尺寸仍是 640×400
    NSGraphicsContext.saveGraphicsState()
    NSGraphicsContext.current = NSGraphicsContext(bitmapImageRep: rep)
    draw()
    NSGraphicsContext.restoreGraphicsState()
    try rep.representation(using: .png, properties: [:])!.write(to: url)
}

let out = URL(fileURLWithPath: CommandLine.arguments.count > 1 ? CommandLine.arguments[1] : ".")
try FileManager.default.createDirectory(at: out, withIntermediateDirectories: true)
try render(scale: 1, to: out.appendingPathComponent("background.png"))
try render(scale: 2, to: out.appendingPathComponent("background@2x.png"))
print("✓ \(out.path)")
