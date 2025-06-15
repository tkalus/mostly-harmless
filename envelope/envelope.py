#!/usr/bin/env python3
"""
Print destination and return addresses on envelopes from a toml input file.

---

Requires:
  - Python 3.11.x or later; tomllib support
  - Pycairo - https://pypi.org/project/pycairo/

  $ pip install pycairo
"""
# /// script
# dependencies = [
#   "pycairo",
# ]
# ///

from __future__ import annotations

import collections.abc as abc
import sys
import tomllib

import cairo

ENVELOPE_SIZES: abc.Mapping[str, str] = {
    "#5": "5.5x3.125",  # 5-1/2 x 3-1/8
    "#6.25": "6x3.5",  # 6 x 3-1/2
    "#6.75": "6.5x3.625",  # 6-1/2 x 3-5/8
    "#7": "6.75x3.75",  # 6-3/4 x 3-3/4
    "#7.75": "7.5x3.93750",  # 7-1/2 x 3-15/16
    "#8.625": "8.625x3.625",  # 8-5/8 x 3-5/8
    "#9": "8.875x3.875",  # 8-7/8 x 3-7/8
    "#10": "9.5x4.125",  # 9-1/2 x 4-1/8
    "#11": "10.375x4.5",  # 10-3/8 x 4-1/2
    "#12": "11x4.75",  # 11 x 4-3/4
    "#14": "11.5x5",  # 11-1/2 x 5
}
FONT_FACE: str = "mono"  # serif and sans-serif also possible
FONT_SIZE: int = 14


def write_envelope_pdf(
    io_out: abc.BinaryIO,
    from_addr: abc.Sequence[str],
    to_addr: abc.Sequence[str],
    width: float,
    height: float,
    font_face: str = FONT_FACE,
    font_size: int = FONT_SIZE,
) -> None:
    """
    Write an envelope-PDF to passed io_out.

    Generate an envlope-PDF of the specified size in float-inches, the to_addr
    centered, the from_addr in the top-left, using the specified font_face and
    font_size (defaults to 14, mono-spaced).
    """
    INCHES_TO_POINTS = 72
    MARGIN = 0.25

    pt_width = width * INCHES_TO_POINTS
    pt_height = height * INCHES_TO_POINTS
    pt_margin = MARGIN * INCHES_TO_POINTS

    surface = cairo.PDFSurface(io_out, pt_width, pt_height)
    ctx = cairo.Context(surface)
    ctx.select_font_face(font_face)
    ctx.set_font_size(font_size)

    # Position the from/return address at top-left
    x = pt_margin
    y = (pt_margin * 1.5) + font_size
    for line in from_addr:
        ctx.move_to(x, y)
        ctx.show_text(line)
        y += font_size

    # Vertically-center the to address ~40% from the left
    x = pt_width * 0.38
    y = (pt_height / 2) - (font_size * (len(to_addr) - 1) / 2)
    for line in to_addr:
        ctx.move_to(x, y)
        ctx.show_text(line)
        y += font_size

    ctx.show_page()
    surface.flush()
    surface.finish()


def _main() -> None:
    in_data = tomllib.loads("""
size = "#10"  # 9-1/2 x 4-1/8 inches
from = [
  "E. L. Brown",
  "1640 Riverside Drive",
  "Hill Valley, CA 91103",
]
to = [
  "Burton Richter",
  "c/o SLAC National Laboratory",
  "2575 Sand Hill Rd",
  "Menlo Park, CA 94025",
]
""")
    if in_file := dict(enumerate(sys.argv)).get(1, ""):
        in_file = "/dev/stdin" if in_file=="-" else in_file
        # Passed a pathname
        with open(in_file, "rb") as f:
            in_data = tomllib.load(f)

    # validate envelope data.
    size = in_data.get("size", "#10")
    from_addr = in_data.get("from", [])
    to_addr = in_data.get("to", [])
    if not all(
        (
            isinstance(size, str),
            isinstance(from_addr, abc.Sequence),
            isinstance(to_addr, abc.Sequence),
            from_addr,
            to_addr,
        )
    ):
        raise TypeError("input data malformed")

    # Parse dimensions, translating common, known envelope sizes.
    size_str = size
    size_str = ENVELOPE_SIZES.get(size_str, size_str)
    width, height = size_str.split(sep="x", maxsplit=2)

    with open("envelope.pdf", "wb") as f:
        write_envelope_pdf(
            io_out=f,
            from_addr=from_addr,
            to_addr=to_addr,
            width=float(width),
            height=float(height),
        )


if __name__ == "__main__":
    _main()
