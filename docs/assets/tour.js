// Responsive, silent player for the deterministic tour artwork.
(function () {
  'use strict';
  const player = document.querySelector('.tour-player');
  if (!player) return;
  const labels = JSON.parse(player.dataset.labels);
  const viewport = player.querySelector('.tour-viewport');
  const film = document.getElementById('film');
  const toggle = player.querySelector('.tour-toggle');
  const slider = player.querySelector('.tour-progress');
  const timer = player.querySelector('.tour-time');
  const chapters = player.querySelector('.tour-chapters');
  const fullscreen = player.querySelector('.tour-fullscreen');
  const status = player.querySelector('.tour-status');
  const reduced = matchMedia('(prefers-reduced-motion: reduce)');
  const params = new URLSearchParams(location.search);
  let duration = 84.5;
  let time = 0;
  let playing = false;
  let frame = 0;
  let previous = 0;
  let initialized = false;
  let poster = true;
  let chapterIndex = -1;

  function layout() {
    const compact = matchMedia('(max-width: 600px)').matches && !params.has('render');
    film.classList.toggle('compact', compact);
    const width = compact ? 960 : 1920;
    const height = compact ? 1200 : 1080;
    film.style.width = width + 'px';
    film.style.height = height + 'px';
    viewport.style.aspectRatio = width + ' / ' + height;
    film.style.transform = 'scale(' + viewport.clientWidth / width + ')';
    if (initialized) {
      window.tourMeasure();
      draw();
    }
  }

  function format(value) {
    const seconds = Math.floor(value);
    return Math.floor(seconds / 60) + ':' + String(seconds % 60).padStart(2, '0');
  }

  function draw() {
    // Keep the closing card visible when playback ends or reaches the last frame.
    window.seek(Math.min(time, duration - 0.4));
    slider.value = time;
    timer.textContent = format(time) + ' / ' + format(Math.ceil(duration));
    slider.setAttribute('aria-valuetext', format(time) + ' / ' + format(Math.ceil(duration)));
    toggle.textContent = playing ? labels.pause : time >= duration ? labels.replay : labels.play;
    const index = window.tourScenes.findIndex(s => time < s.start + s.dur);
    const current = index < 0 ? window.tourScenes.length - 1 : index;
    chapters.value = current;
    if (chapterIndex !== current) {
      chapterIndex = current;
      status.textContent = chapters.options[current].textContent;
    }
  }

  function pause() {
    playing = false;
    cancelAnimationFrame(frame);
    frame = 0;
    draw();
  }

  function tick(now) {
    if (!playing) return;
    time = Math.min(duration, time + (now - previous) / 1000);
    previous = now;
    if (time >= duration) return pause();
    draw();
    frame = requestAnimationFrame(tick);
  }

  function play() {
    if (playing || !initialized) return;
    if (poster || time >= duration) time = 0;
    poster = false;
    playing = true;
    previous = performance.now();
    draw();
    frame = requestAnimationFrame(tick);
  }

  function move(value) {
    poster = false;
    time = Math.max(0, Math.min(duration, value));
    previous = performance.now();
    draw();
  }

  toggle.addEventListener('click', () => playing ? pause() : play());
  slider.addEventListener('input', () => move(+slider.value));
  chapters.addEventListener('change', () => move(window.tourScenes[+chapters.value].start + 0.5));
  player.addEventListener('keydown', event => {
    if (/^(INPUT|SELECT|BUTTON|A)$/.test(event.target.tagName)) return;
    if (event.code === 'Space') {
      event.preventDefault();
      playing ? pause() : play();
    } else if (event.key === 'ArrowLeft' || event.key === 'ArrowRight') {
      event.preventDefault();
      move(time + (event.key === 'ArrowRight' ? 5 : -5));
    }
  });
  player.tabIndex = 0;
  document.addEventListener('keydown', event => {
    if (event.key === 'Escape' && params.has('embed') && parent !== window && !document.fullscreenElement) {
      event.preventDefault();
      parent.postMessage('tour:close', location.origin);
    }
  });

  if (!player.requestFullscreen || !document.fullscreenEnabled) fullscreen.hidden = true;
  fullscreen.addEventListener('click', async () => {
    try {
      if (document.fullscreenElement) await document.exitFullscreen();
      else await player.requestFullscreen();
    } catch (_) { /* Keep the player usable when full screen is unavailable. */ }
  });
  document.addEventListener('fullscreenchange', () => {
    const label = document.fullscreenElement ? labels.exit_fullscreen : labels.fullscreen;
    fullscreen.setAttribute('aria-label', label);
    fullscreen.title = label;
    layout();
  });
  document.addEventListener('visibilitychange', () => {
    if (document.hidden && initialized) pause();
  });
  // Closing the parent dialog pauses its iframe without reloading the artwork.
  window.addEventListener('message', event => {
    if (event.origin !== location.origin || event.source !== parent || !initialized) return;
    if (event.data === 'tour:pause') pause();
    if (event.data === 'tour:play' && !reduced.matches) play();
  });
  reduced.addEventListener('change', () => { if (reduced.matches && initialized) pause(); });

  if (params.has('render')) document.documentElement.classList.add('rendering');
  layout();
  if ('ResizeObserver' in window) new ResizeObserver(layout).observe(viewport);
  else window.addEventListener('resize', layout);
  window.tourReady = window.ready().then(info => {
    duration = info.duration;
    slider.max = duration;
    initialized = true;
    time = window.tourScenes[1].start + 3.4;
    draw();
    player.querySelector('.tour-controls').hidden = params.has('render');
    // The parent decides whether its dialog is still open. A visitor may have
    // closed it while images were loading; that must never start hidden playback.
    if (params.has('embed') && parent !== window) parent.postMessage('tour:ready', location.origin);
    return info;
  });
})();
