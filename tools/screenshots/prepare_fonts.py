"""从本机的 macOS 系统字体生成截图渲染用的静态字体文件。

    python3 tools/screenshots/prepare_fonts.py <输出目录>

Flutter 的测试环境只认单个字体文件里的第一个字形集，也不会把字重映射到可变字体的轴上，
直接加载系统字体的话标题会是"假粗体"。这里把 PingFang.ttc 拆成单独的字重，
把可变字体 SFNS.ttf 实例化成固定字重，让截图与真实 App 的渲染一致。

这些是 Apple 的系统字体，只在本机生成到临时目录，不要提交进仓库。
需要 fonttools：pip install fonttools
"""

import glob
import os
import sys

from fontTools.ttLib import TTCollection, TTFont
from fontTools.varLib.instancer import instantiateVariableFont

out = sys.argv[1] if len(sys.argv) > 1 else "fonts"
os.makedirs(out, exist_ok=True)

# ── 中文：PingFang SC 的三个字重 ──
pingfang = (glob.glob("/System/Library/AssetsV2/com_apple_MobileAsset_Font*/*/AssetData/PingFang.ttc")
            + glob.glob("/System/Library/Fonts/PingFang.ttc"))[0]
wanted = {"PingFang SC Regular": "PingFang-Regular", "PingFang SC Medium": "PingFang-Medium",
          "PingFang SC Semibold": "PingFang-Semibold"}
for font in TTCollection(pingfang).fonts:
    name = font["name"].getDebugName(4)
    if name in wanted:
        font.save(os.path.join(out, wanted.pop(name) + ".ttf"))
if wanted:
    sys.exit(f"PingFang.ttc 里缺少：{', '.join(wanted)}")

# ── 日文：冬青黑体（Hiragino Sans）的三个字重，日文界面的截图用它，汉字是日文字形 ──
for weight, label in [("W3", "Regular"), ("W5", "Medium"), ("W6", "Semibold")]:
    for font in TTCollection(f"/System/Library/Fonts/ヒラギノ角ゴシック {weight}.ttc").fonts:
        if font["name"].getDebugName(4) == f"Hiragino Sans {weight}":
            font.save(os.path.join(out, f"Hiragino-{label}.ttf"))
            break
    else:
        sys.exit(f"缺少 Hiragino Sans {weight}")

# ── 等宽：Menlo，用于配对码与安全指纹 ──
menlo = {"Menlo Regular": "Menlo-Regular", "Menlo Bold": "Menlo-Bold"}
for font in TTCollection("/System/Library/Fonts/Menlo.ttc").fonts:
    name = font["name"].getDebugName(4)
    if name in menlo:
        font.save(os.path.join(out, menlo.pop(name) + ".ttf"))

# ── 西文：SF Pro 实例化为固定字重（opsz 取正文尺寸）──
for weight, label in [(400, "Regular"), (500, "Medium"), (600, "Semibold")]:
    vf = TTFont("/System/Library/Fonts/SFNS.ttf")
    axes = {a.axisTag: a for a in vf["fvar"].axes}
    loc = {"wght": weight}
    if "opsz" in axes:
        loc["opsz"] = max(axes["opsz"].minValue, min(17, axes["opsz"].maxValue))
    for tag, a in axes.items():
        loc.setdefault(tag, a.defaultValue)
    instantiateVariableFont(vf, loc, inplace=True)
    vf.save(os.path.join(out, f"SF-{label}.ttf"))

print("✓", out)
