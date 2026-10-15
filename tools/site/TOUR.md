# Website tour

Each language's homepage opens a localized, silent 85-second tour. The same
tour is available at `tour/` under each language root. The homepage keeps its
short copy-and-paste animation and places **Watch the full tour** next to the
download button. The full tour's artwork and screenshots load only on click.

- `tour.html` holds the shared scene markup with `{{semantic_key}}` placeholders.
- `content/tour.py` holds five complete translations and the player labels.
  Shared product names and comparison text come from the existing site copy.
- `docs/assets/tour.css` contains the artwork, locale fonts, player controls,
  and a separate portrait composition for phones.
- `docs/assets/tour-scenes.js` renders each frame as a function of time.
- `docs/assets/tour.js` handles playback, seeking, chapters, scaling, full screen,
  reduced motion, and the parent dialog's playback messages.

Run `python3 tools/site/build.py` after changing copy or markup, and commit the
source and generated pages together. The original files in `tools/promo/` are
independent of the website and are not required to build it.

For frame rendering, open a generated tour URL with `?render`, await
`window.tourReady`, and call `window.seek(seconds)`. This fixes the artwork at
1920 × 1080 and hides the page header and player controls. `window.ready()`
returns the duration and original audio cue timings for an external renderer.

Check all nine chapters in every language at desktop and phone widths. Also
check lazy loading, pause/resume, seeking, chapter selection, full screen,
closing before screenshots finish loading, closing/reopening the dialog,
reduced motion, and keyboard focus returning to the homepage link.
