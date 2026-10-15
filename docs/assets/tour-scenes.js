// Deterministic artwork: window.seek(seconds) also supports frame rendering.

'use strict';

// ── 工具 ──
const $ = (s) => document.querySelector(s);
const $$ = (s) => [...document.querySelectorAll(s)];
const clamp = (x, a = 0, b = 1) => Math.min(b, Math.max(a, x));
const lerp = (a, b, k) => a + (b - a) * k;
const linear = (x) => x;
const easeOut = (x) => 1 - Math.pow(1 - x, 3);
const easeInOut = (x) => (x < 0.5 ? 4 * x * x * x : 1 - Math.pow(-2 * x + 2, 3) / 2);
const easeBack = (x) => 1 + 2.2 * Math.pow(x - 1, 3) + 1.2 * Math.pow(x - 1, 2);
// t 从 a 走到 b 时，0→1
const p = (t, a, b, e = easeOut) => e(clamp((t - a) / (b - a)));

// 淡入并上移；k = 0 时在下方 dy 处、透明
function rise(el, k, dy = 30) {
  el.style.opacity = k;
  el.style.transform = `translateY(${(1 - k) * dy}px)`;
}

function mixColor(a, b, k) {
  const pa = a.match(/\w\w/g).map((h) => parseInt(h, 16));
  const pb = b.match(/\w\w/g).map((h) => parseInt(h, 16));
  return `rgb(${pa.map((v, i) => Math.round(lerp(v, pb[i], k))).join(',')})`;
}

// 按键：0.32 秒内按下再弹起，按下时变蓝
function press(el, t, at) {
  const x = (t - at) / 0.32;
  let k = 0;
  if (x > 0 && x < 1) k = x < 0.4 ? x / 0.4 : 1 - (x - 0.4) / 0.6;
  el.style.transform = `translateY(${k * 2}px) scale(${1 - 0.08 * k})`;
  el.style.background = mixColor('3b3d43', '3d8bea', k);
}

// ── 音效提示点：audio.py 按这些时间点合成音效 ──
const CUES = [];
function cue(scene, t, type, extra = {}) {
  CUES.push({ t: +(scene.start + t).toFixed(3), type, ...extra });
}

// ── 图标：与 assets/brand/icon.svg 同一套几何，前面那张蓝纸可以单独动 ──
const LOGO = `
  <defs>
    <linearGradient id="LID-body" x1="0" y1="100" x2="0" y2="924" gradientUnits="userSpaceOnUse">
      <stop offset="0" stop-color="#FBFBFA"/><stop offset="1" stop-color="#E8E8E6"/>
    </linearGradient>
    <filter id="LID-bs" x="-20%" y="-20%" width="140%" height="140%"><feDropShadow dx="0" dy="12" stdDeviation="14" flood-color="#000" flood-opacity="0.28"/></filter>
    <filter id="LID-ss" x="-20%" y="-20%" width="140%" height="140%"><feDropShadow dx="0" dy="10" stdDeviation="12" flood-color="#000" flood-opacity="0.22"/></filter>
  </defs>
  <g class="lg-body"><rect x="100" y="100" width="824" height="824" rx="184" fill="url(#LID-body)" filter="url(#LID-bs)"/></g>
  <g class="lg-back">
    <rect x="230" y="220" width="340" height="420" rx="44" fill="#F5F5F4" stroke="#2C2C2E" stroke-width="28"/>
    <g fill="#2C2C2E" fill-opacity="0.85"><rect x="290" y="310" width="200" height="30" rx="15"/><rect x="290" y="390" width="150" height="30" rx="15"/><rect x="290" y="470" width="104" height="30" rx="15"/></g>
  </g>
  <g class="lg-front">
    <rect x="454" y="384" width="340" height="420" rx="44" fill="#0A66D8" filter="url(#LID-ss)"/>
    <g fill="#FFFFFF" fill-opacity="0.95"><rect x="514" y="474" width="200" height="30" rx="15"/><rect x="514" y="554" width="150" height="30" rx="15"/><rect x="514" y="634" width="104" height="30" rx="15"/></g>
  </g>`;
for (const svg of $$('.logo')) svg.innerHTML = LOGO.replaceAll('LID', svg.id);

// 图标出场：底板放大淡入，后面的纸出现，蓝纸从它身上「复制」出来滑到右下
function logoIn(svg, t, t0) {
  const body = svg.querySelector('.lg-body');
  const back = svg.querySelector('.lg-back');
  const front = svg.querySelector('.lg-front');
  const kb = p(t, t0, t0 + 0.7);
  body.style.opacity = kb;
  body.setAttribute('transform', `translate(512 512) scale(${0.9 + 0.1 * kb}) translate(-512 -512)`);
  back.style.opacity = p(t, t0 + 0.3, t0 + 0.8);
  front.style.opacity = p(t, t0 + 0.7, t0 + 0.9, linear);
  const m = p(t, t0 + 0.75, t0 + 1.45, easeBack);
  front.setAttribute('transform', `translate(${(1 - m) * -224} ${(1 - m) * -164})`);
}

// ── 镜头 ──
const SCENES = [];
let clock = 0;
function scene(id, dur, render) {
  const s = { id, el: document.getElementById(id), start: clock, dur, render };
  clock += dur;
  SCENES.push(s);
  return s;
}

// 1 痛点
const S1 = scene('s1', 10.2, (t) => {
  const head = $('#s1-head');
  const mv = p(t, 2.8, 3.6, easeInOut);
  head.style.transform = `translateY(${-290 * mv}px) scale(${lerp(1, 0.72, mv)})`;
  head.style.opacity = 1 - p(t, 7.9, 8.4);
  rise($('#s1-l1'), p(t, 0.3, 1.0));
  rise($('#s1-l2'), p(t, 1.1, 1.8));
  $$('#s1 .pain').forEach((card, i) => {
    const k = p(t, 3.4 + i * 0.35, 4.0 + i * 0.35);
    const s = p(t, 6.3 + i * 0.28, 6.75 + i * 0.28, easeInOut);
    card.querySelector('i').style.transform = `scaleX(${s})`;
    card.style.opacity = k * lerp(1, 0.45, s) * (1 - p(t, 7.9, 8.4));
    card.style.transform = `translateY(${(1 - k) * 30}px)`;
  });
  rise($('#s1-l3'), p(t, 8.6, 9.3));
});
[0, 1, 2, 3].forEach((i) => cue(S1, 3.4 + i * 0.35, 'tick'));
[0, 1, 2, 3].forEach((i) => cue(S1, 6.3 + i * 0.28, 'strike'));

// 2 品牌
const S2 = scene('s2', 5.2, (t) => {
  logoIn($('#s2-logo'), t, 0.2);
  rise($('#s2-name'), p(t, 1.4, 2.0));
  rise($('#s2-tag'), p(t, 2.0, 2.6));
  rise($('#s2-os'), p(t, 2.6, 3.2));
});
cue(S2, 0.0, 'music');
cue(S2, 1.2, 'pop');

// 3 演示：三个回合，复制的东西化作一张蓝纸沿弧线飞过去
const ROUNDS = {
  a: { src: '#a-src .doc', dst: '#a-dst .doc', fly: 2.35, dir: 1, avoid: [] },
  b: { src: '#b-src', dst: '#b-dst', fly: 8.5, dir: -1, avoid: ['#w2 .appwin', '#m2 .appwin'] },
  c: { src: '#c-src', dst: '#c-dst', fly: 13.7, dir: 1, avoid: ['#m3 .appwin', '#w3 .appwin'] },
};
const FLY_DUR = 1.5;
const LIFT = 200;

function center(sel) {
  const r = $(sel).getBoundingClientRect();
  const f = $('#film').getBoundingClientRect();
  const scale = f.width / $('#film').offsetWidth;
  return { x: (r.left - f.left + r.width / 2) / scale, y: (r.top - f.top + r.height / 2) / scale };
}

function bez(S, C, T, u) {
  return {
    x: (1 - u) * (1 - u) * S.x + 2 * u * (1 - u) * C.x + u * u * T.x,
    y: (1 - u) * (1 - u) * S.y + 2 * u * (1 - u) * C.y + u * u * T.y,
  };
}

function rect(sel) {
  const r = $(sel).getBoundingClientRect();
  const f = $('#film').getBoundingClientRect();
  const scale = f.width / $('#film').offsetWidth;
  return { l: (r.left - f.left) / scale, t: (r.top - f.top) / scale, r: (r.right - f.left) / scale, b: (r.bottom - f.top) / scale };
}

// 量一次起点、终点，在弧线上每隔 26px 放一个点；落在 App 窗口上的点不画，免得压住文字
function measure() {
  const stage = $('#stage');
  const destination = $('#a-dst');
  const paste = $('#b-paste');
  const saved = [stage.style.transform, destination.style.transform, paste.style.display];
  stage.style.transform = 'none';
  destination.style.transform = 'translate(-50%, -50%)';
  paste.style.display = '';
  const trail = $('#trail');
  trail.replaceChildren();
  for (const [key, r] of Object.entries(ROUNDS)) {
    const avoid = r.avoid.map(rect);
    const covered = (q) => avoid.some((a) => q.x > a.l - 14 && q.x < a.r + 14 && q.y > a.t - 14 && q.y < a.b + 14);
    r.S = center(r.src);
    r.T = center(r.dst);
    r.C = $('#film').classList.contains('compact')
      ? { x: r.S.x + 240 * r.dir, y: (r.S.y + r.T.y) / 2 }
      : { x: (r.S.x + r.T.x) / 2, y: Math.min(r.S.y, r.T.y) - 2 * LIFT };
    r.apex = bez(r.S, r.C, r.T, 0.5);
    const g = document.createElementNS('http://www.w3.org/2000/svg', 'g');
    g.id = `trail-${key}`;
    let last = r.S;
    let run = 0;
    for (let i = 1; i <= 400; i++) {
      const u = i / 400;
      const q = bez(r.S, r.C, r.T, u);
      run += Math.hypot(q.x - last.x, q.y - last.y);
      last = q;
      // 两头各留一段，不压在文件上
      if (run >= 26 && u > 0.08 && u < 0.92 && !covered(q)) {
        run = 0;
        const c = document.createElementNS('http://www.w3.org/2000/svg', 'circle');
        c.setAttribute('cx', q.x.toFixed(1));
        c.setAttribute('cy', q.y.toFixed(1));
        c.setAttribute('r', '3');
        c.dataset.u = u;
        g.appendChild(c);
      }
    }
    trail.appendChild(g);
    r.g = g;
  }
  [stage.style.transform, destination.style.transform, paste.style.display] = saved;
}

// 第 t 秒时，某个回合的纸飞到哪、轨迹露出多少；vis 是这个回合整体的可见度
function flight(r, t, vis) {
  const raw = (t - r.fly) / FLY_DUR;
  const u = easeInOut(clamp(raw));
  for (const c of r.g.children) {
    c.style.opacity = raw > 0 && u >= +c.dataset.u ? 0.5 * vis : 0;
  }
  if (raw <= 0 || raw >= 1) return false;
  const q = bez(r.S, r.C, r.T, u);
  const scale = raw < 0.5 ? 0.55 + 0.45 * (raw / 0.5) : 1 - 0.4 * ((raw - 0.5) / 0.5);
  const rot = 8 * Math.sin(Math.PI * raw) * r.dir;
  const fly = $('#fly');
  fly.style.opacity = raw < 0.12 ? raw / 0.12 : raw > 0.88 ? (1 - raw) / 0.12 : 1;
  fly.style.transform = `translate(${q.x - 40}px, ${q.y - 51}px) rotate(${rot}deg) scale(${scale})`;
  return true;
}

const S3 = scene('s3', 17.3, (t) => {
  const stage = $('#stage');
  stage.style.transform = `translateY(${(1 - p(t, 0.1, 0.9)) * 40}px)`;

  // 回合 A：Mac 上的文件 → Windows
  const va = 1 - p(t, 5.9, 6.3);
  $('#m1').style.opacity = va;
  $('#w1').style.opacity = va;
  const sel = p(t, 1.5, 1.7) * (1 - p(t, 4.4, 4.8));
  $('#a-name').style.background = `rgba(61,139,234,${sel})`;
  const keysA = $$('#m1 kbd');
  press(keysA[0], t, 1.9);
  press(keysA[1], t, 2.02);
  const keysA2 = $$('#w1 kbd');
  press(keysA2[0], t, 3.95);
  press(keysA2[1], t, 4.07);
  const land = p(t, 4.3, 4.9);
  const dstA = $('#a-dst');
  dstA.style.opacity = land;
  dstA.style.transform = `translate(-50%, ${lerp(-40, -50, land)}%) scale(${lerp(0.85, 1, land)})`;

  // 回合 B：Windows 记事本里的文字 → Mac 备忘录
  const vb = p(t, 6.4, 6.9) * (1 - p(t, 11.6, 12.0));
  $('#m2').style.opacity = vb;
  $('#w2').style.opacity = vb;
  $('#b-sel1').style.backgroundSize = `${p(t, 7.3, 7.6, linear) * 100}% 100%`;
  $('#b-sel2').style.backgroundSize = `${p(t, 7.6, 7.9, linear) * 100}% 100%`;
  const keysB = $$('#w2 kbd');
  press(keysB[0], t, 8.1);
  press(keysB[1], t, 8.22);
  const keysB2 = $$('#m2 kbd');
  press(keysB2[0], t, 10.05);
  press(keysB2[1], t, 10.17);
  const pasted = p(t, 10.4, 10.9);
  $('#b-paste').style.opacity = pasted;
  $('#b-paste').style.display = t < 10.4 ? 'none' : '';
  $('#b-caret').style.opacity = Math.floor(t * 2) % 2 === 0 ? 1 : 0;

  // 回合 C：Mac 上的截图 → Windows 里的 Word
  const vc = p(t, 12.1, 12.6);
  $('#m3').style.opacity = vc;
  $('#w3').style.opacity = vc;
  $('#c-ring').style.opacity = p(t, 13.0, 13.2) * (1 - p(t, 15.6, 16.0));
  const keysC = $$('#m3 kbd');
  press(keysC[0], t, 13.3);
  press(keysC[1], t, 13.42);
  const keysC2 = $$('#w3 kbd');
  press(keysC2[0], t, 15.25);
  press(keysC2[1], t, 15.37);
  const img = p(t, 15.6, 16.2);
  $('#c-dst').style.opacity = img;
  $('#c-dst').style.transform = `scale(${lerp(0.9, 1, img)})`;

  // 纸和轨迹
  const flying = flight(ROUNDS.a, t, va) | flight(ROUNDS.b, t, vb) | flight(ROUNDS.c, t, vc);
  if (!flying) $('#fly').style.opacity = 0;
  const label = $('#arc-label');
  const cur = t < 6.2 ? ROUNDS.a : t < 12.0 ? ROUNDS.b : ROUNDS.c;
  const lv = t < 6.2 ? p(t, 2.9, 3.4) * va : t < 12.0 ? p(t, 9.0, 9.5) * vb : p(t, 14.2, 14.7) * vc;
  label.style.left = `${cur.apex.x}px`;
  label.style.top = `${cur.apex.y - 64}px`;
  label.style.opacity = lv;

  rise($('#cap-a'), p(t, 0.9, 1.5) * (1 - p(t, 5.6, 6.0)), 20);
  rise($('#cap-b'), p(t, 6.6, 7.2) * (1 - p(t, 11.5, 11.9)), 20);
  rise($('#cap-c'), p(t, 12.3, 12.9), 20);
});
// 按键音：Mac 在左声道，Windows 在右声道
[[1.9, 2.02, 3.95, 4.07, -1], [8.1, 8.22, 10.05, 10.17, 1], [13.3, 13.42, 15.25, 15.37, -1]].forEach(([a, b, c, d, side]) => {
  for (const k of [a, b]) cue(S3, k, 'key', { pan: 0.5 * side });
  for (const k of [c, d]) cue(S3, k, 'key', { pan: -0.5 * side });
});
for (const r of Object.values(ROUNDS)) cue(S3, r.fly, 'whoosh', { dur: FLY_DUR, dir: r.dir });
cue(S3, 4.3, 'land', { pan: 0.5 });
cue(S3, 10.4, 'land', { pan: -0.5 });
cue(S3, 15.6, 'land', { pan: 0.5 });

// 4 能同步什么
const S4 = scene('s4', 6.8, (t) => {
  rise($('#s4-h'), p(t, 0.3, 0.9));
  $$('#s4 .tile').forEach((el, i) => rise(el, p(t, 0.9 + i * 0.16, 1.5 + i * 0.16), 40));
  rise($('#s4-sub'), p(t, 2.4, 3.0));
});
[0, 1, 2, 3, 4].forEach((i) => cue(S4, 0.9 + i * 0.16, 'tick'));

// 5 复制记录
const S5 = scene('s5', 8.6, (t, dur) => {
  rise($('#s5-text h2'), p(t, 0.3, 0.9));
  $$('#s5 .points li').forEach((el, i) => rise(el, p(t, 1.5 + i * 0.5, 2.1 + i * 0.5), 20));
  const shot = $('#s5-shot');
  const k = p(t, 0.5, 1.5);
  shot.style.opacity = k;
  shot.style.transform = `translateX(${(1 - k) * 120}px) scale(${1 + 0.03 * (t / dur)})`;
  $('#s5-mark').style.opacity = p(t, 3.4, 3.8);
});
[0, 1, 2, 3].forEach((i) => cue(S5, 1.5 + i * 0.5, 'tick'));

// 6 配对
const S6 = scene('s6', 9.2, (t) => {
  rise($('#s6-h'), p(t, 0.3, 0.9));
  rise($('#s6-sub1'), p(t, 0.6, 1.2) * (1 - p(t, 4.3, 4.6)), 16);
  rise($('#s6-sub2'), p(t, 4.6, 5.2), 16);
  const shot = $('#s6-shot');
  const k = p(t, 0.6, 1.6);
  shot.style.opacity = k * lerp(1, 0.35, p(t, 4.6, 5.2));
  shot.style.transform = `translateY(${(1 - k) * 50}px)`;
  const fade = 1 - p(t, 4.4, 4.8);
  $$('#s6 .pulse').forEach((c, i) => {
    if (t < 1.2) { c.setAttribute('r', 0); c.style.opacity = 0; return; }
    const ph = ((t - 1.2) / 2.4 + i / 3) % 1;
    c.setAttribute('r', lerp(30, 225, ph).toFixed(1));
    c.style.opacity = (1 - ph) * 0.55 * fade * p(t, 1.2, 1.6);
  });
  const card = $('#s6-card');
  const kc = p(t, 4.6, 5.4);
  card.style.opacity = kc;
  card.style.transform = `translateY(${(1 - kc) * 60}px)`;
});
cue(S6, 4.6, 'pop');

// 7 连接方式
const S7 = scene('s7', 9.6, (t) => {
  rise($('#s7-h'), p(t, 0.3, 0.9));
  $$('#s7 .way').forEach((el, i) => {
    const k = p(t, 1.0 + i * 0.35, 1.6 + i * 0.35);
    rise(el, k, 40);
    const dot = el.querySelector('.p-dot');
    const ph = (((t - 1.6 - i * 0.35) / 1.4) % 1 + 1) % 1;
    dot.setAttribute('cx', lerp(+dot.dataset.x1, +dot.dataset.x2, easeInOut(ph)).toFixed(2));
    dot.setAttribute('cy', dot.dataset.y);
    dot.style.opacity = k * Math.sin(Math.PI * ph);
  });
  rise($('#s7-foot'), p(t, 4.0, 4.7));
});
[0, 1, 2].forEach((i) => cue(S7, 1.0 + i * 0.35, 'tick'));

// 8 对比
const S8 = scene('s8', 9.0, (t) => {
  rise($('#s8-h'), p(t, 0.3, 0.9));
  const rows = $$('#s8 .row');
  rows.forEach((el, i) => rise(el, p(t, 0.8 + i * 0.3, 1.3 + i * 0.3), 16));
  $('#s8-col').style.opacity = p(t, 1.0, 1.8);
});
[1, 2, 3, 4, 5, 6, 7].forEach((i) => cue(S8, 0.8 + i * 0.3, 'tick', { gain: 0.6 }));

// 9 结尾
const S9 = scene('s9', 8.6, (t) => {
  logoIn($('#s9-logo'), t, 0.2);
  rise($('#s9-name'), p(t, 1.1, 1.7));
  rise($('#s9-tag'), p(t, 1.6, 2.2));
  rise($('#s9-chips'), p(t, 2.2, 2.8));
  rise($('#s9-links'), p(t, 2.9, 3.5));
  rise($('#s9-ask'), p(t, 4.2, 4.8));
});
cue(S9, 0.0, 'finale');
cue(S9, 1.2, 'pop');

const DURATION = clock;
const FADE = 0.35;

function seek(T) {
  for (const s of SCENES) {
    const t = T - s.start;
    if (t < 0 || t > s.dur) {
      s.el.style.visibility = 'hidden';
      continue;
    }
    s.el.style.visibility = 'visible';
    s.el.style.opacity = Math.min(p(t, 0, FADE, linear), 1 - p(t, s.dur - FADE, s.dur, linear));
    s.render(t, s.dur);
  }
}

async function ready() {
  await document.fonts.ready;
  await Promise.all($$('img').map((im) => im.decode().catch(() => {})));
  measure();
  return { duration: DURATION, cues: CUES };
}

window.seek = seek;
window.ready = ready;


window.tourScenes = SCENES.map(({ id, start, dur }) => ({ id, start, dur }));
window.tourMeasure = measure;
