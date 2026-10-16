# PasteFerry 2.1 portrait launch film

A 25.5-second, 1080×1920, 30 fps video with Chinese narration, burned-in captions,
a matching SRT, a portrait cover and original synthesized music. No external
media or model service is used. The computer screens are labeled motion demos,
not recordings of a user's clipboard. All information uses the public website;
no administrative hostname, test identity or credential is included.

The opening immediately shows a cross-network copy/paste result. The next cuts
show the official-service switch, direct-first routing, the current relay
allowance and the official download address. Text stays above the bottom
platform overlay area and clear of the right action rail. Do not claim unlimited
relay or that private production storage is proven by the public source.

Requirements: macOS with the Tingting voice (`say`), Chrome, Node.js,
`puppeteer-core`, Python 3 with NumPy, and FFmpeg with libx264. Set `NODE_PATH` to
the puppeteer installation, and optionally `NODE_BIN`, `PYTHON_BIN`, `CHROME`.

```bash
./tools/shorts/build.sh /absolute/output/path
```

Use `node tools/shorts/render.js /output/path stills` to inspect the key frames
before rendering. `story.json` owns the narration, timings and captions. Visuals
in `video.html` use the same deterministic timeline. Refresh both when editing.
Build outputs stay outside the 14 release attachments.
