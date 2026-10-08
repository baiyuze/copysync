// 生成网站的分享图（Open Graph，1200×630）：左侧图标与标题，右侧露出一部分真实截图。
//
//   swift tools/site/og.swift
//
// 输出 docs/assets/og.png（中文）与 docs/assets/og-en.png（英文）。

import AppKit

let W: CGFloat = 1200, H: CGFloat = 630
let paper = NSColor(srgbRed: 0xF5 / 255.0, green: 0xF5 / 255.0, blue: 0xF4 / 255.0, alpha: 1)
let ink = NSColor(srgbRed: 0x1D / 255.0, green: 0x1D / 255.0, blue: 0x1F / 255.0, alpha: 1)
let dim = NSColor(srgbRed: 0x6E / 255.0, green: 0x6E / 255.0, blue: 0x73 / 255.0, alpha: 1)

func render(title: String, note: String, to path: String) throws {
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: Int(W), pixelsHigh: Int(H),
                               bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false,
                               colorSpaceName: .deviceRGB, bytesPerRow: 0, bitsPerPixel: 0)!
    NSGraphicsContext.saveGraphicsState()
    NSGraphicsContext.current = NSGraphicsContext(bitmapImageRep: rep)
    defer { NSGraphicsContext.restoreGraphicsState() }

    paper.setFill()
    NSRect(x: 0, y: 0, width: W, height: H).fill()

    // 右侧截图：带窗口外框的那张，向右下出血
    if let shot = NSImage(contentsOfFile: "docs/assets/screenshots/history.png") {
        let h: CGFloat = 560
        let w = h * shot.size.width / shot.size.height
        shot.draw(in: NSRect(x: 600, y: H - 70 - h, width: w, height: h))
    }

    // 左侧：图标、产品名、一句话、注脚（坐标按左上原点换算）
    func y(_ top: CGFloat, _ height: CGFloat) -> CGFloat { H - top - height }
    NSImage(contentsOfFile: "assets/brand/icon-1024.png")?
        .draw(in: NSRect(x: 64, y: y(72, 112), width: 112, height: 112))

    let name = NSAttributedString(string: "CopySync", attributes: [
        .font: NSFont.systemFont(ofSize: 34, weight: .semibold), .foregroundColor: ink,
    ])
    name.draw(at: NSPoint(x: 80, y: y(200, 42)))

    let para = NSMutableParagraphStyle()
    para.lineSpacing = 6
    let headline = NSAttributedString(string: title, attributes: [
        .font: NSFont.systemFont(ofSize: 54, weight: .bold), .foregroundColor: ink,
        .paragraphStyle: para, .kern: -1,
    ])
    headline.draw(in: NSRect(x: 80, y: y(262, 150), width: 540, height: 150))

    let foot = NSAttributedString(string: note, attributes: [
        .font: NSFont.systemFont(ofSize: 22, weight: .regular), .foregroundColor: dim,
    ])
    foot.draw(at: NSPoint(x: 80, y: y(470, 28)))

    try rep.representation(using: .png, properties: [:])!.write(to: URL(fileURLWithPath: path))
}

try render(title: "在这台 Mac 复制，\n到那台 Mac 粘贴。", note: "免费开源，macOS 13 及以上",
           to: "docs/assets/og.png")
try render(title: "Copy on this Mac.\nPaste on that one.", note: "Free and open source, macOS 13+",
           to: "docs/assets/og-en.png")
print("✓ docs/assets/og.png, og-en.png")
