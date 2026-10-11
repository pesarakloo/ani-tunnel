#!/usr/bin/env python3
"""Rasterize ani tunnel's SVG mark. Requires cairosvg for asset development only."""
from pathlib import Path
import sys
import cairosvg
root = Path(__file__).resolve().parents[4]
output = Path(sys.argv[1]) if len(sys.argv) > 1 else Path(__file__).parent
svg = (root / "img/ani-mark.svg").read_bytes()
for name, size in [("icon-192.png", 192), ("icon-512.png", 512), ("apple-touch-icon.png", 180), ("icon-maskable-512.png", 512)]:
    content = svg
    if "maskable" in name:
        inner = svg.decode().split(">", 1)[1].rsplit("</svg>", 1)[0]
        content = ('<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 128 128"><rect width="128" height="128" fill="#061020"/><g transform="translate(18 18) scale(.72)">' + inner + '</g></svg>').encode()
    cairosvg.svg2png(bytestring=content, write_to=str(output/name), output_width=size, output_height=size)
