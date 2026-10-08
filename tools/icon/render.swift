// 绘制 CopySync 的 App 图标，输出各尺寸 PNG。
//
//   swift tools/icon/render.swift <输出目录>
//
// 造型：两张错位的纸。后一张描边、前一张实色，表示「这台复制、那台出现」。
// 网站用的 tools/icon/icon.svg 与这里的几何参数一一对应，改动时两边同步。
// 画布按 Apple 的 macOS 图标网格：1024 画布，主体 824×824 居中，四周留出阴影空间。

import AppKit

let graphite = NSColor(srgbRed: 0x2C / 255.0, green: 0x2C / 255.0, blue: 0x2E / 255.0, alpha: 1)
let blue = NSColor(srgbRed: 0x0A / 255.0, green: 0x66 / 255.0, blue: 0xD8 / 255.0, alpha: 1)
let bodyTop = NSColor(srgbRed: 0xFB / 255.0, green: 0xFB / 255.0, blue: 0xFA / 255.0, alpha: 1)
let bodyBottom = NSColor(srgbRed: 0xE8 / 255.0, green: 0xE8 / 255.0, blue: 0xE6 / 255.0, alpha: 1)
let paper = NSColor(srgbRed: 0xF5 / 255.0, green: 0xF5 / 255.0, blue: 0xF4 / 255.0, alpha: 1)

/// 以 1024 画布、左上角为原点的坐标描述图形，绘制时再翻转到 CoreGraphics 的左下原点。
func draw(in ctx: CGContext) {
    ctx.translateBy(x: 0, y: 1024)
    ctx.scaleBy(x: 1, y: -1)

    // ── 主体：圆角方块 + 投影 ──
    let body = CGPath(roundedRect: CGRect(x: 100, y: 100, width: 824, height: 824),
                      cornerWidth: 184, cornerHeight: 184, transform: nil)
    ctx.saveGState()
    // 阴影偏移不受上面的翻转影响，仍按左下原点计，负值才是向下
    ctx.setShadow(offset: CGSize(width: 0, height: -12), blur: 28,
                  color: NSColor(white: 0, alpha: 0.28).cgColor)
    ctx.addPath(body)
    ctx.setFillColor(bodyBottom.cgColor)
    ctx.fillPath()
    ctx.restoreGState()

    // 自上而下极轻的明暗，模拟 macOS 图标统一的顶光
    ctx.saveGState()
    ctx.addPath(body)
    ctx.clip()
    let grad = CGGradient(colorsSpace: CGColorSpace(name: CGColorSpace.sRGB),
                          colors: [bodyTop.cgColor, bodyBottom.cgColor] as CFArray,
                          locations: [0, 1])!
    ctx.drawLinearGradient(grad, start: CGPoint(x: 512, y: 100), end: CGPoint(x: 512, y: 924),
                           options: [])
    ctx.restoreGState()

    // ── 后一张纸：描边 ──
    // 两张纸沿对角线拉开，只在一角重叠：两边的三行字都完整露出，
    // 读作「同一份内容在两处」，而不是系统里那个「拷贝」符号
    let back = CGRect(x: 230, y: 220, width: 340, height: 420)
    let backPath = CGPath(roundedRect: back, cornerWidth: 44, cornerHeight: 44, transform: nil)
    ctx.addPath(backPath)
    ctx.setFillColor(paper.cgColor)
    ctx.fillPath()
    ctx.addPath(backPath)
    ctx.setStrokeColor(graphite.cgColor)
    ctx.setLineWidth(28)
    ctx.strokePath()
    lines(in: ctx, sheet: back, color: graphite.withAlphaComponent(0.85))

    // ── 前一张纸：实色，压在后一张上 ──
    let front = CGRect(x: 454, y: 384, width: 340, height: 420)
    let frontPath = CGPath(roundedRect: front, cornerWidth: 44, cornerHeight: 44, transform: nil)
    ctx.saveGState()
    ctx.setShadow(offset: CGSize(width: 0, height: -10), blur: 24,
                  color: NSColor(white: 0, alpha: 0.22).cgColor)
    ctx.addPath(frontPath)
    ctx.setFillColor(blue.cgColor)
    ctx.fillPath()
    ctx.restoreGState()
    lines(in: ctx, sheet: front, color: NSColor(white: 1, alpha: 0.95))
}

/// 纸上的三行字：长、中、短。
func lines(in ctx: CGContext, sheet: CGRect, color: NSColor) {
    let x = sheet.minX + 60
    let widths: [CGFloat] = [200, 150, 104]
    ctx.setFillColor(color.cgColor)
    for (i, w) in widths.enumerated() {
        let y = sheet.minY + 90 + CGFloat(i) * 80
        let bar = CGPath(roundedRect: CGRect(x: x, y: y, width: w, height: 30),
                         cornerWidth: 15, cornerHeight: 15, transform: nil)
        ctx.addPath(bar)
        ctx.fillPath()
    }
}

func render(size: Int, to url: URL) throws {
    let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: size, pixelsHigh: size,
                               bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true,
                               isPlanar: false, colorSpaceName: .deviceRGB,
                               bytesPerRow: 0, bitsPerPixel: 0)!
    let gctx = NSGraphicsContext(bitmapImageRep: rep)!
    let ctx = gctx.cgContext
    ctx.interpolationQuality = .high
    ctx.setShouldAntialias(true)
    ctx.scaleBy(x: CGFloat(size) / 1024, y: CGFloat(size) / 1024)
    draw(in: ctx)
    gctx.flushGraphics()
    guard let png = rep.representation(using: .png, properties: [:]) else {
        throw NSError(domain: "icon", code: 1)
    }
    try png.write(to: url)
}

let out = URL(fileURLWithPath: CommandLine.arguments.count > 1 ? CommandLine.arguments[1] : ".")
try FileManager.default.createDirectory(at: out, withIntermediateDirectories: true)
for size in [16, 32, 64, 128, 180, 256, 512, 1024] {
    try render(size: size, to: out.appendingPathComponent("icon_\(size).png"))
}
print("✓ \(out.path)")
