// 首屏动画、代码块的「复制」按钮、分页（部署教程按系统分页、快速开始按部署方式分页）、语言菜单、记住选过的语言。
// 没有脚本时页面照样完整：动画停在复制完成的样子，分页的内容依次排开。
(function () {
  // Load the localized tour only when requested. The href remains a normal
  // standalone page for browsers without dialog support or JavaScript.
  var tourLink = document.querySelector('[data-tour-open]');
  var tourDialog = document.querySelector('.tour-dialog');
  if (tourLink && tourDialog && typeof tourDialog.showModal === 'function') {
    var tourFrame = tourDialog.querySelector('iframe');
    tourLink.addEventListener('click', function (event) {
      if (event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey || event.altKey) return;
      event.preventDefault();
      tourDialog.showModal();
      document.documentElement.classList.add('tour-is-open');
      if (!tourFrame.hasAttribute('src')) tourFrame.src = tourFrame.dataset.src;
      else tourFrame.contentWindow.postMessage('tour:play', location.origin);
    });
    tourDialog.addEventListener('click', function (event) {
      if (event.target === tourDialog) tourDialog.close();
    });
    tourDialog.addEventListener('close', function () {
      tourFrame.contentWindow.postMessage('tour:pause', location.origin);
      document.documentElement.classList.remove('tour-is-open');
      tourLink.focus();
    });
    window.addEventListener('message', function (event) {
      if (event.origin !== location.origin || event.source !== tourFrame.contentWindow) return;
      if (event.data === 'tour:close') tourDialog.close();
      if (event.data === 'tour:ready') tourFrame.contentWindow.postMessage(tourDialog.open ? 'tour:play' : 'tour:pause', location.origin);
    });
  }

  var stage = document.querySelector("[data-play]");

  // 按两台设备实际的位置算出纸飞行的起点、终点与弧线
  function layout() {
    var src = stage.querySelector(".mac .doc");
    var dst = stage.querySelector(".win .doc");
    var sheet = stage.querySelector(".sheet-fly");
    var path = stage.querySelector(".arc path");
    var label = stage.querySelector(".arc-label");
    var box = stage.getBoundingClientRect();
    var a = src.getBoundingClientRect();
    var b = dst.getBoundingClientRect();
    var w = sheet.offsetWidth;
    var h = sheet.offsetHeight;
    var sx = a.left + a.width / 2 - box.left;
    var sy = a.top + a.height / 2 - box.top;
    var tx = b.left + b.width / 2 - box.left;
    var ty = b.top + b.height / 2 - box.top;
    var lift = Math.max(40, (tx - sx) * 0.28);
    stage.style.setProperty("--sx", sx - w / 2 + "px");
    stage.style.setProperty("--sy", sy - h / 2 + "px");
    stage.style.setProperty("--tx", tx - w / 2 + "px");
    stage.style.setProperty("--ty", ty - h / 2 + "px");
    stage.style.setProperty("--lift", lift + "px");
    var peak = Math.min(sy, ty) - lift * 1.6;
    path.setAttribute("d", "M" + sx + " " + (sy - a.height / 2) +
      " Q" + (sx + tx) / 2 + " " + peak + " " + tx + " " + (ty - b.height / 2));
    label.style.setProperty("--ly", peak + lift * 0.55 - 8 + "px");
  }

  function play() {
    layout();
    stage.classList.remove("play");
    void stage.offsetWidth; // 强制重排，让动画从头再来
    stage.classList.add("play");
  }

  if (stage) {
    layout();
    window.addEventListener("resize", layout);
    // 滚动到可见时才开始，免得还没看到就播完了
    if ("IntersectionObserver" in window) {
      var seen = new IntersectionObserver(function (entries) {
        if (entries[0].isIntersecting) {
          play();
          seen.disconnect();
        }
      }, { threshold: 0.35 });
      seen.observe(stage);
    } else {
      play();
    }
    var replay = document.querySelector(".replay");
    if (replay) replay.addEventListener("click", play);
  }

  document.querySelectorAll(".code").forEach(function (block) {
    var button = block.querySelector(".copy");
    var pre = block.querySelector("pre");
    if (!button || !pre || !navigator.clipboard) {
      if (button) button.hidden = true;
      return;
    }
    var label = button.textContent;
    button.addEventListener("click", function () {
      navigator.clipboard.writeText(pre.innerText.trim()).then(function () {
        button.textContent = button.dataset.done || "Copied";
        setTimeout(function () { button.textContent = label; }, 1600);
      });
    });
  });

  // 分页。链接里带了 #windows、#plan-lan 之类就打开那一页；否则按系统分页（data-tabs="os"）的按访问者的系统，
  // 按部署方式分页（data-tabs="plan"）的打开第一页
  document.querySelectorAll("[data-tabs]").forEach(function (box) {
    var list = box.querySelector("[role=tablist]");
    var tabs = Array.prototype.slice.call(list.querySelectorAll("[role=tab]"));
    var panels = tabs.map(function (t) { return document.getElementById(t.getAttribute("aria-controls")); });

    function select(i, focus) {
      tabs.forEach(function (t, j) {
        t.setAttribute("aria-selected", i === j ? "true" : "false");
        t.tabIndex = i === j ? 0 : -1;
        panels[j].hidden = i !== j;
      });
      if (focus) tabs[i].focus();
    }

    tabs.forEach(function (t, i) {
      t.addEventListener("click", function () {
        select(i);
        history.replaceState(null, "", "#" + panels[i].id);
      });
      t.addEventListener("keydown", function (e) {
        var step = e.key === "ArrowRight" ? 1 : e.key === "ArrowLeft" ? -1 : 0;
        if (step) {
          e.preventDefault();
          select((i + step + tabs.length) % tabs.length, true);
        }
      });
    });

    var ids = panels.map(function (p) { return p.id; });
    var pick = ids.indexOf(location.hash.slice(1));
    var linked = pick >= 0;
    if (!linked && box.dataset.tabs !== "os") {
      pick = 0;
    } else if (!linked) {
      var ua = navigator.userAgent;
      pick = /Windows/.test(ua) ? 0 : /Macintosh/.test(ua) ? 1 : /Linux|X11/.test(ua) && !/Android/.test(ua) ? 2 : 0;
    }
    list.hidden = false;
    box.classList.add("tabs-on");
    select(pick);
    // 其他系统的内容收起来了，页面变短，重新对准链接指向的位置
    if (linked) box.scrollIntoView();
  });

  // 下载页：猜这台 Mac 是哪种芯片，把对应的按钮放在前面、改成主按钮。
  // 浏览器不直接告诉芯片，只能看显卡型号：Chrome、Firefox 会给出 Intel 或 Apple M…，
  // Safari 只说 Apple GPU，猜不出来就保持默认（Apple 芯片在前）
  var chips = document.querySelector("[data-mac-chip]");
  if (chips && /Macintosh/.test(navigator.userAgent)) {
    try {
      var gl = document.createElement("canvas").getContext("webgl");
      var info = gl && gl.getExtension("WEBGL_debug_renderer_info");
      var gpu = info ? String(gl.getParameter(info.UNMASKED_RENDERER_WEBGL)) : "";
      if (/Intel|AMD|Radeon/i.test(gpu) && !/Apple M/i.test(gpu)) {
        var buttons = chips.querySelectorAll(".button");
        buttons[0].classList.add("quiet");
        buttons[1].classList.remove("quiet");
        chips.insertBefore(buttons[1], buttons[0]);
      }
    } catch (e) {}
  }

  // 语言菜单：点菜单外面或按 Esc 收起
  var langs = document.querySelector("details.lang");
  if (langs) {
    document.addEventListener("click", function (e) {
      if (langs.open && !langs.contains(e.target)) langs.open = false;
    });
    langs.addEventListener("keydown", function (e) {
      if (e.key === "Escape" && langs.open) {
        langs.open = false;
        langs.querySelector("summary").focus();
      }
    });
  }

  // 选过语言，以后打开首页就直接是这种语言，不再按浏览器语言跳转
  document.querySelectorAll("a[data-lang]").forEach(function (link) {
    link.addEventListener("click", function () {
      try { localStorage.setItem("pasteferry-lang", link.dataset.lang); } catch (e) {}
    });
  });
})();
