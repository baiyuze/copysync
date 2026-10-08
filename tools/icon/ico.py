"""把 render.swift --windows 输出的各尺寸 PNG 打包成 Windows 的 .ico。

    python3 tools/icon/ico.py <PNG 所在目录> <输出的 .ico>

每个尺寸直接以 PNG 存进 ico（Windows Vista 起支持），保留透明度，也不需要额外的库。
"""

import re
import struct
import sys
from pathlib import Path


def main(src: Path, dst: Path) -> None:
    images = []
    for p in sorted(src.glob("windows_*.png"), key=lambda p: int(re.findall(r"\d+", p.name)[0])):
        data = p.read_bytes()
        # PNG 的 IHDR 里是宽高
        w, h = struct.unpack(">II", data[16:24])
        images.append((w, h, data))
    if not images:
        sys.exit(f"{src} 里没有 windows_*.png")

    header = struct.pack("<HHH", 0, 1, len(images))  # 保留、类型 1 = 图标、数量
    offset = len(header) + 16 * len(images)
    entries, blobs = b"", b""
    for w, h, data in images:
        # 宽高为 256 时按约定写 0
        entries += struct.pack("<BBBBHHII", w % 256, h % 256, 0, 0, 1, 32, len(data), offset)
        blobs += data
        offset += len(data)
    dst.write_bytes(header + entries + blobs)
    print(f"✓ {dst}（{', '.join(str(w) for w, _, _ in images)}）")


if __name__ == "__main__":
    main(Path(sys.argv[1]), Path(sys.argv[2]))
