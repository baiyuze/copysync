// 首屏演示与代码块的「复制」按钮。没有它们页面也完整可读。
(function () {
  var demo = document.querySelector(".demo");
  var replay = document.querySelector(".replay");

  function play() {
    if (!demo) return;
    demo.classList.remove("play");
    void demo.offsetWidth; // 强制重排，让动画从头再来
    demo.classList.add("play");
  }

  if (demo) {
    // 滚动到可见时才开始，避免用户还没看到就播完了
    if ("IntersectionObserver" in window) {
      var seen = new IntersectionObserver(function (entries) {
        if (entries[0].isIntersecting) {
          play();
          seen.disconnect();
        }
      }, { threshold: 0.5 });
      seen.observe(demo);
    } else {
      play();
    }
  }
  if (replay) replay.addEventListener("click", play);

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
})();
