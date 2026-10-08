// 给渲染出的界面截图加上 macOS 窗口外框：圆角、窗口阴影、红绿灯按钮与细边框。
//
//   swift tools/screenshots/frame.swift <输入目录> <输出目录>
//
// 输入是 2 倍图（960×640 的窗口渲染为 1920×1280）。红绿灯的位置量自真实窗口：
// 直径 12pt，圆心距窗口左上角 14pt，间隔 20pt。

import AppKit

let scale: CGFloat = 2
let radius: CGFloat = 10 * scale
let margin = (left: 48 * scale, right: 48 * scale, top: 32 * scale, bottom: 64 * scale)

func rgb(_ hex: UInt32, _ a: CGFloat = 1) -> CGColor {
    CGColor(srgbRed: CGFloat((hex >> 16) & 0xFF) / 255, green: CGFloat((hex >> 8) & 0xFF) / 255,
            blue: CGFloat(hex & 0xFF) / 255, alpha: a)
}

func frame(_ input: URL, to output: URL) throws {
    guard let src = NSBitmapImageRep(data: try Data(contentsOf: input))?.cgImage else { return }
    let dark = input.lastPathComponent.contains("dark")
    let w = CGFloat(src.width), h = CGFloat(src.height)
    let W = Int(w + margin.left + margin.right), H = Int(h + margin.top + margin.bottom)

    let ctx = CGContext(data: nil, width: W, height: H, bitsPerComponent: 8, bytesPerRow: 0,
                        space: CGColorSpace(name: CGColorSpace.sRGB)!,
                        bitmapInfo: CGImageAlphaInfo.premultipliedLast.rawValue)!
    // CoreGraphics 原点在左下；窗口矩形按左上原点算好再换算
    let win = CGRect(x: margin.left, y: margin.bottom, width: w, height: h)
    let shape = CGPath(roundedRect: win, cornerWidth: radius, cornerHeight: radius, transform: nil)

    // macOS 的窗口阴影：大而淡，略向下
    ctx.saveGState()
    ctx.setShadow(offset: CGSize(width: 0, height: -18 * scale), blur: 40 * scale,
                  color: rgb(0x000000, dark ? 0.45 : 0.22))
    ctx.addPath(shape)
    ctx.setFillColor(rgb(0x000000))
    ctx.fillPath()
    ctx.restoreGState()

    ctx.saveGState()
    ctx.addPath(shape)
    ctx.clip()
    ctx.draw(src, in: win)
    ctx.restoreGState()

    // 细边框：浅色窗口是一圈淡灰，深色窗口是一圈淡白
    ctx.addPath(CGPath(roundedRect: win.insetBy(dx: 0.5, dy: 0.5), cornerWidth: radius,
                       cornerHeight: radius, transform: nil))
    ctx.setStrokeColor(dark ? rgb(0xFFFFFF, 0.14) : rgb(0x000000, 0.16))
    ctx.setLineWidth(1)
    ctx.strokePath()

    // 红绿灯
    let lights: [(UInt32, UInt32)] = [(0xFF5F57, 0xE2463F), (0xFEBC2E, 0xE1A116), (0x28C840, 0x14AE2B)]
    for (i, (fill, stroke)) in lights.enumerated() {
        let cx = win.minX + (14 + CGFloat(i) * 20) * scale
        let cy = win.maxY - 14 * scale
        let r = 6 * scale
        let dot = CGRect(x: cx - r, y: cy - r, width: r * 2, height: r * 2)
        ctx.setFillColor(rgb(fill))
        ctx.fillEllipse(in: dot)
        ctx.setStrokeColor(rgb(stroke, 0.9))
        ctx.setLineWidth(1)
        ctx.strokeEllipse(in: dot.insetBy(dx: 0.5, dy: 0.5))
    }

    let rep = NSBitmapImageRep(cgImage: ctx.makeImage()!)
    try rep.representation(using: .png, properties: [:])!.write(to: output)
}

let args = CommandLine.arguments
let inDir = URL(fileURLWithPath: args[1]), outDir = URL(fileURLWithPath: args[2])
try FileManager.default.createDirectory(at: outDir, withIntermediateDirectories: true)
for name in try FileManager.default.contentsOfDirectory(atPath: inDir.path) where name.hasSuffix(".png") {
    let src = inDir.appendingPathComponent(name), dst = outDir.appendingPathComponent(name)
    // card- 开头的是单独渲染的对话框，自带圆角与阴影，不套窗口外框
    if name.hasPrefix("card-") {
        try? FileManager.default.removeItem(at: dst)
        try FileManager.default.copyItem(at: src, to: dst)
    } else {
        try frame(src, to: dst)
    }
}
print("✓ \(outDir.path)")
