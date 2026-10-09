// 首屏动画、代码块的「复制」按钮、记住选过的语言。没有脚本时页面照样完整：动画停在复制完成的样子。
(function () {
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

  // 点过语言链接，以后打开首页就直接是这种语言，不再按浏览器语言跳转
  document.querySelectorAll("a[data-lang]").forEach(function (link) {
    link.addEventListener("click", function () {
      try { localStorage.setItem("copysync-lang", link.dataset.lang); } catch (e) {}
    });
  });
})();
