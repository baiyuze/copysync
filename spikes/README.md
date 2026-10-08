# M-1 Spike：验证 macOS 剪贴板的两个致命假设

CopySync 的「纯懒模式」（复制时不传输，对端按下粘贴键才回源拉取）在技术上
**只能**通过剪贴板的延迟渲染（promise）机制实现——操作系统不提供「剪贴板被读取」
的通知，promise 回调是唯一能感知「对方按下了 Cmd+V」的时机。

这两个 spike 验证该路径是否真的走得通。**若假设不成立，整个架构需要重新设计。**

本机环境：macOS 15.7.3 (Sequoia) · Apple clang 17.0.0 · Xcode 26.x

---

## Spike A — 剪贴板读取的隐私提示（`privacy/`）

**要验证什么**：macOS 15 Sequoia 起，App 以编程方式读取*其他* App 放入剪贴板的内容时，
系统会弹出「"某某" 想要从其他 App 粘贴」的确认框。CopySync 的 daemon 需要持续监听剪贴板，
如果每次复制都弹窗，产品直接不可用。

**分阶段探测**，逐级判断哪一步触发提示：

| 阶段 | 操作 | 预期 |
|---|---|---|
| 1 | 读 `changeCount` | 应**不**触发（只是个计数器） |
| 2 | 读 `pasteboard.types` | 应**不**触发（只是类型列表） |
| 3 | `detectPatternsForPatterns:` (macOS 14+) | 官方宣称**不**触发，用于"窥探但不读取" |
| 4 | 读 `public.file-url` 的实际内容 | **很可能触发** ← 这是我们真正需要的 |

**成败判据**：
- 阶段 1/2 不触发 → 至少「感知到剪贴板变了」是免费的
- 阶段 4 触发但选「每次都允许」后永久静默 → 可接受，首次引导用户授权一次
- 阶段 4 每次都弹且无法永久授权 → **架构需调整**（改为用户主动触发同步，而非自动监听）

## Spike B — `NSFilePromiseProvider` 的实际可用性（`promise/`）

**要验证什么**：把 file promise 放进剪贴板后，目标 App 粘贴时能否正确触发回调、
回调里能否**异步延迟数秒**（模拟跨网络拉取）而不被系统超时杀掉。

程序把一个 promise 写进剪贴板，delegate 回调里先 `sleep 5s` 再写文件，
然后你去各个 App 里 `Cmd+V`，观察终端日志和 App 行为。

**成败判据**：
- Finder 能粘出文件、显示进度条、不超时 → 核心路径成立
- 5 秒延迟被杀 → 需要缩短首字节时间，或改为预落盘
- 大量常用 App 不支持 promise → 需要对这些 App 降级为「先落盘再给路径」

**兼容性矩阵**（在 `results.md` 里填写实测结果）：
Finder / 微信 / 飞书 / Mail / VS Code / 钉钉 / Safari 上传框 / 终端拖拽区

---

## 运行

```bash
./build.sh
```

编译产物是两个 `.app` bundle（**不是**裸命令行工具）。这点很重要：macOS 的 TCC 隐私机制
以 bundle identifier 为授权主体，裸二进制的提示行为与真实 App 不一致，测出来的结论不可用。

```bash
./build/SpikeA.app/Contents/MacOS/SpikeA    # 前台运行看日志
./build/SpikeB.app/Contents/MacOS/SpikeB
```

两个程序都在前台打印日志，`Ctrl+C` 退出。
