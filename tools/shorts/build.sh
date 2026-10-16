#!/usr/bin/env bash
# Reproducible 1080x1920, 30 fps launch film. All output stays outside release/.
set -euo pipefail
cd "$(dirname "$0")/../.."
SHORTS_OUT=${1:-dist/shorts-2.1}
NODE_BIN=${NODE_BIN:-node}
PYTHON_BIN=${PYTHON_BIN:-python3}
mkdir -p "$SHORTS_OUT"
"$NODE_BIN" tools/shorts/render.js "$SHORTS_OUT" stills
"$PYTHON_BIN" tools/shorts/audio.py "$SHORTS_OUT"
"$NODE_BIN" tools/shorts/render.js "$SHORTS_OUT/frames" frames
ffmpeg -y -hide_banner -loglevel warning -framerate 30 -i "$SHORTS_OUT/frames/frame-%05d.jpg" -i "$SHORTS_OUT/mix.wav" \
  -vf 'scale=out_color_matrix=bt709:out_range=tv,format=yuv420p' \
  -c:v libx264 -preset medium -crf 18 -g 60 -colorspace bt709 -color_primaries bt709 -color_trc bt709 \
  -af 'loudnorm=I=-16:TP=-1.5:LRA=7' -c:a aac -b:a 192k -ar 48000 -t 25.5 -movflags +faststart \
  "$SHORTS_OUT/PasteFerry-2.1-douyin-1080x1920.mp4"
