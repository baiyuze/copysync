# M-1 Spike 实测结果

环境：macOS 15.7.3 (Sequoia) · Apple clang 17.0.0 · 2026-09-21

---

## Spike A — 剪贴板隐私机制

### 背景更正（比计划中的假设更严重）

计划里写的是「Sequoia 起读剪贴板会弹提示」，实测后发现真实情况是：

- macOS **15.4** 引入该机制，但目前是 **developer preview，默认关闭**
- 需 `defaults write <bundle_id> EnablePasteboardPrivacyDeveloperPreview -bool yes` 才生效
- **下一个 macOS 大版本将强制执行**

所以今天不做任何事也能正常工作，但这是颗定时炸弹。CopySync 的 daemon 正是
「后台程序化读剪贴板」，完全命中被拦截的场景，必须现在就按新机制设计。

参考：[macOS 16 will warn when apps access pasteboard](https://appletreats.substack.com/p/macos-16-will-warn-when-apps-access)

### 实测数据（隐私机制开启时）

写方是独立 bundle id 的真实 GUI App（SpikeC），排除了 pbcopy/osascript
不被 TCC 视作「其他 App」的干扰。

| 阶段 | 操作 | 耗时 | 是否触发授权 |
|---|---|---|---|
| 0 | `pb.accessBehavior` | 0.001s | 否 ✓ |
| 1 | `pb.changeCount` | 0.000s | 否 ✓ |
| 2 | `pb.types`（类型列表） | 0.003s | **否 ✓** |
| 3 | `detectMetadataForTypes:` | 0.021s | **否 ✓** |
| 4 | `readObjectsForClasses:`（读 file-url） | **阻塞 90s 被杀** | **是** |

关闭隐私机制时，阶段 4 仅耗时 0.001s —— 差异完全来自授权拦截。

### 结论 1：类型分类完全免费

`detectMetadataForTypes:` 成功返回 `ContentType = public.plain-text`，
并能正确判断 `conformsToType:` 的文件夹 / 图片 / 纯文本归属。

意味着 CopySync 可以**免提示地**完成：
- 感知剪贴板变化（changeCount）
- 判断复制的是文件 / 文件夹 / 文本 / 图片（detectMetadata + types）

只有拿**源文件真实路径**时才需要授权 —— 而这是 A 端回源所必需的。

### 结论 2：必须在主线程调用剪贴板 API ★

隐私机制开启后，在后台 dispatch queue 上访问剪贴板会**直接崩溃**：

```
reportAccessBehavior → AppKit → ViewBridge -[NSCFRunLoopSemaphore wait]
                     → NSException raise → EXC_BAD_ACCESS (SIGSEGV)
```

原因：剪贴板 API 要经 ViewBridge 与系统 UI 服务通信（以便显示授权框），
其内部 `NSCFRunLoopSemaphore` 依赖 run loop，裸 dispatch queue 上无 run loop 可用。

改为主线程 NSTimer 轮询后崩溃消失，变为正常阻塞等待授权。

**对 Go daemon 的直接约束**：剪贴板 cgo 调用不能在任意 goroutine 里发起，
必须调度到锁定的主线程（`runtime.LockOSThread` + 主 run loop），
或通过 `dispatch_async(dispatch_get_main_queue(), …)` 弹回主线程。

同理，等待异步回调时不能裸阻塞，要 spin run loop：

```objc
while (dispatch_semaphore_wait(sem, DISPATCH_TIME_NOW) != 0) {
    [NSRunLoop.currentRunLoop runMode:NSDefaultRunLoopMode
                           beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.02]];
}
```

否则一旦 completionHandler 投递到主队列即死锁，且授权弹窗无法显示。

### 结论 3：授权状态会自动迁移

实测 `accessBehavior` 变化：`0 (Default)` → 触发一次弹窗 → `1 (Ask)`。

与 SDK 文档一致：首次触发后 App 开始出现在
「系统设置 → 隐私与安全性 → 从其他 App 粘贴」，用户可切换
Ask / AlwaysAllow / AlwaysDeny。

**设计要点**：读取前先查 `accessBehavior`，若不是 `AlwaysAllow` 就不要贸然读取
（否则会无限阻塞在弹窗上），而应在 UI 里引导用户去系统设置授权。

### 结论 4：粘贴场景天然豁免

SDK 文档明确：`.ask` 和 `.alwaysDeny` 下，
「user originated and paste related」的访问**永远被允许且不通知**。

意味着 B 端 promise 回调里的剪贴板访问不受限制 —— 只有 A 端的后台监听需要授权。

---

---

## Spike B — 延迟渲染（★ 结论：纯懒模式不成立）

### 测了两条路径，都无法感知「用户按下 Cmd+V」

**路径 1：`NSFilePromiseProvider`**

写入剪贴板后类型为 `Apple files promise pasteboard type` 等 promise 私有类型。
Finder 粘贴**毫无反应**，连 `fileNameForType` 都不触发 —— Finder 不认为
剪贴板里有可粘贴的文件。该 API 实际只服务于**拖放**，不覆盖剪贴板粘贴路径。

（用自建的 `NSFilePromiseReceiver` 消费方倒是能触发 `fileNameForType`，
但文件传输始终不完成，最终超时。）

**路径 2：`declareTypes:owner:` + `provideDataForType:`**

声明 Finder 认识的 `public.file-url` / `NSFilenamesPboardType`，不提供数据。
`provideDataForType:` **确实被触发了**，且回调里阻塞 5 秒未被系统杀掉。

但关键在触发时机——实测：

| 阶段 | 操作 | 回调次数 |
|---|---|---|
| 1 | 只写剪贴板，静置 12 秒，不做任何操作 | **1 次，在 0.0 秒** |
| 2 | 打开 Finder 窗口，再等 8 秒 | 仍 1 次，无新增 |
| 3 | 用户真正按 Cmd+V | **不再回调**，直接用缓存 |

用户粘贴出的文件内容自证：`落盘时刻 21:35:51 / 粘贴延迟 0.1 秒`，
而 Finder 窗口 21:35:49 才打开 —— 数据在用户按键前就备好了。

**系统在写入剪贴板后 0.1 秒内立即兑现 promise 并缓存，不等粘贴。**

### 对架构的影响

操作系统不提供「剪贴板被读取」的通知，延迟渲染是唯一能感知粘贴动作的机制。
它失效 ⇒ **无法知道对方何时按下 Cmd+V** ⇒ 「对方不粘贴就不传输」不成立。

可用的部分：`declareTypes:owner:` 能把文件正确送进 Finder 粘贴。
丢掉的是「延迟」，不是「粘贴」。

### 未测（因已不影响结论）

- 拖放场景下 `NSFilePromiseProvider` 是否真正延迟到「放下」时刻
  —— 拖放不是本项目要的交互
- 是否为 Universal Clipboard 导致的立即读取
  —— 即便是，也不能要求用户关闭该系统功能

---

---

## 后续在实现中发现的坑

### pasteboard 服务会被 promise owner 搞成僵尸状态 ★

M-1 的 SpikeB 用 `declareTypes:owner:` 声明了剪贴板所有权，进程被 kill 后，
系统 pasteboard 留下了「有类型声明但无数据」的空壳：

```
$ osascript -e 'clipboard info'
«class utf8», 0, «class ut16», 2, string, 0, Unicode text, 0
                ↑ 所有类型长度都是 0
$ pbpaste          # 什么都没有
```

这个状态下**任何程序都读不到剪贴板**，包括 `pbpaste` 和 SpikeA 自己。
排查时一度误判为自己的 cgo 代码有问题，实际上：

```bash
killall pboard    # 重启 pasteboard 服务，再复制一次即恢复
```

教训：调试剪贴板问题时，先用 `pbpaste` 确认剪贴板本身是否正常，
再怀疑自己的代码。这条已写入 README 的排障章节。

### 传输层的三个真实故障

均在 macOS 实机联调中撞出，各自对应回归测试
（`internal/transport/p2p/stream_test.go`）：

| 现象 | 根因 |
|---|---|
| 控制消息完全收不到 | `SettingEngine.DetachDataChannels()` 是全局开关，开启后所有通道的 `OnMessage`/`Send` 失效 |
| 发 202 字节、收 0 字节 | pion 的 `OnMessage` 不缓冲，handler 设置前到达的消息直接丢弃 |
| zstd 报 "compressed size too big" | 「先缓冲后回放」的实现里，回放期间新消息插队，顺序被打乱 |
| 20 MB 传到 786 KB 卡死 | `OnBufferedAmountLow` 边沿触发，缓冲提前回落导致丢失唤醒 |

### 配对完成瞬间的信令竞态

两端确认配对必然有先后（实测相差 47 ms），先确认方发出的 offer 会被
后确认方当作「未配对设备的消息」丢弃，握手随即石沉大海，只能等 30 秒
超时重来。解法是暂存配对确认前到达的信令、确认后重放——配对完成到
建立直连的实测耗时从 30 秒降到 56 毫秒。

---

## 待确认

- [ ] Windows 侧 `CFSTR_FILECONTENTS` 虚拟文件是否也有同样的「立即兑现」问题
      （Windows 的 `IStream` 延迟渲染机制与 macOS 不同，可能仍然成立）
