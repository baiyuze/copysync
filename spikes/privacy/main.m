// Spike A — macOS 剪贴板隐私机制探测
//
// 背景：macOS 15.4 引入了剪贴板隐私保护，目前是 developer preview（默认关闭），
//      下一个 macOS 大版本将强制执行。CopySync 的 daemon 恰好是
//      「后台程序化读取剪贴板」——完全命中被拦截的场景，必须提前验证。
//
// 要回答的问题：
//   Q1. 哪些操作免提示？（changeCount / types / detectMetadata / 读内容）
//   Q2. detectMetadataForTypes 能否拿到足够的信息做类型分类？
//   Q3. 触发提示后 accessBehavior 如何变化？用户授权后能否永久静默？
//   Q4. 每个阶段耗时多少？（耗时 > 1s 说明弹窗出现并等待了用户交互）
//
// 开启 preview 后再跑一次，即可预演未来 macOS 上的行为：
//   defaults write com.copysync.spike.privacy EnablePasteboardPrivacyDeveloperPreview -bool yes

#import <Cocoa/Cocoa.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>

static const NSTimeInterval kPollInterval = 0.3;
static const NSTimeInterval kBlockThreshold = 1.0;  // 超过此耗时即判定「弹窗出现过」

static void banner(const char *s) {
    fprintf(stderr, "\n════════════════════════════════════════════════════════\n");
    fprintf(stderr, " %s\n", s);
    fprintf(stderr, "════════════════════════════════════════════════════════\n");
    fflush(stderr);
}

static void line(const char *fmt, ...) {
    va_list ap; va_start(ap, fmt);
    fprintf(stderr, "           → ");
    vfprintf(stderr, fmt, ap);
    fprintf(stderr, "\n");
    va_end(ap);
    fflush(stderr);
}

// 打印阶段标题并返回起始时刻，配合 endStage 测量耗时
static NSDate *beginStage(int n, const char *what) {
    fprintf(stderr, "\n  [阶段 %d] %s\n", n, what);
    fflush(stderr);
    return NSDate.date;
}

static void endStage(NSDate *t0) {
    NSTimeInterval dt = -[t0 timeIntervalSinceNow];
    if (dt > kBlockThreshold) {
        fprintf(stderr, "           ⏱  耗时 %.2f 秒 ← 🔔 阻塞了！几乎可以确定弹出了授权框\n", dt);
    } else {
        fprintf(stderr, "           ⏱  耗时 %.3f 秒（无阻塞，未弹窗）\n", dt);
    }
    fflush(stderr);
}

// 在主线程安全地等待异步回调。
//
// 不能用 dispatch_semaphore_wait 直接阻塞主线程：
//   1. 若 completionHandler 恰好投递到主队列，会直接死锁；
//   2. 系统授权弹窗依赖主 run loop 才能显示和响应。
// 正确做法是 spin run loop —— 等待期间继续处理事件。
static BOOL waitSpinningRunLoop(dispatch_semaphore_t sem, NSTimeInterval timeout) {
    NSDate *deadline = [NSDate dateWithTimeIntervalSinceNow:timeout];
    while (dispatch_semaphore_wait(sem, DISPATCH_TIME_NOW) != 0) {
        if (deadline.timeIntervalSinceNow < 0) return NO;
        [NSRunLoop.currentRunLoop runMode:NSDefaultRunLoopMode
                               beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.02]];
    }
    return YES;
}

static const char *behaviorName(NSInteger b) {
    switch (b) {
        case 0: return "Default    （从未触发过提示，不出现在系统设置中）";
        case 1: return "Ask        （每次程序化访问都询问）";
        case 2: return "AlwaysAllow（永久放行，静默）";
        case 3: return "AlwaysDeny （永久拒绝）";
        default: return "未知";
    }
}

static void reportAccessBehavior(NSPasteboard *pb) {
    if (@available(macOS 15.4, *)) {
        line("accessBehavior = %ld  %s", (long)pb.accessBehavior,
             behaviorName(pb.accessBehavior));
    } else {
        line("accessBehavior 不可用（需 macOS 15.4+）");
    }
}

static void probe(NSPasteboard *pb, NSInteger change) {
    banner([[NSString stringWithFormat:@"剪贴板变化 #%ld  %@", (long)change,
             [NSDateFormatter localizedStringFromDate:NSDate.date
                                            dateStyle:NSDateFormatterNoStyle
                                            timeStyle:NSDateFormatterMediumStyle]] UTF8String]);

    // ── 阶段 0：当前授权状态 ───────────────────────────────────
    NSDate *t = beginStage(0, "查询 accessBehavior（当前授权状态）");
    reportAccessBehavior(pb);
    endStage(t);

    // ── 阶段 1：changeCount ───────────────────────────────────
    // 「感知到剪贴板变了」的最低成本操作。若此处弹窗，任何后台监听都不可行。
    t = beginStage(1, "读 changeCount（纯计数器）");
    line("changeCount = %ld", (long)pb.changeCount);
    endStage(t);

    // ── 阶段 2：types ─────────────────────────────────────────
    // 传统的类型判断方式。隐私机制生效后这一步是否还免费，是个关键问题。
    t = beginStage(2, "读 types（类型列表，传统方式）");
    NSArray<NSPasteboardType> *types = pb.types;
    line("%lu 种类型：", (unsigned long)types.count);
    for (NSPasteboardType ty in types) fprintf(stderr, "             · %s\n", ty.UTF8String);
    endStage(t);

    // ── 阶段 3：detectMetadataForTypes ────────────────────────
    // Apple 官方保证「不通知用户」的元数据探测。
    // 文档示例恰好就是 CopySync 的用例：判断文件引用的 content type。
    // 若此处能拿到足够信息，则「分类」环节完全免提示。
    if (@available(macOS 15.4, *)) {
        t = beginStage(3, "detectMetadataForTypes（官方保证不提示）← CopySync 依赖这个");
        dispatch_semaphore_t sem = dispatch_semaphore_create(0);
        [pb detectMetadataForTypes:[NSSet setWithObject:NSPasteboardMetadataTypeContentType]
                 completionHandler:^(NSDictionary<NSPasteboardMetadataType, id> *md, NSError *err) {
            if (err) {
                line("错误：%s", err.localizedDescription.UTF8String);
            } else {
                UTType *ct = md[NSPasteboardMetadataTypeContentType];
                line("ContentType = %s", ct ? ct.identifier.UTF8String : "(nil)");
                if (ct) {
                    line("  是文件夹? %s   是图片? %s   是纯文本? %s",
                         [ct conformsToType:UTTypeFolder] ? "是" : "否",
                         [ct conformsToType:UTTypeImage]  ? "是" : "否",
                         [ct conformsToType:UTTypePlainText] ? "是" : "否");
                }
            }
            dispatch_semaphore_signal(sem);
        }];
        if (!waitSpinningRunLoop(sem, 30)) line("⚠️ 等待超时（30s）");
        endStage(t);
    }

    // ── 阶段 4：读取实际内容 ───────────────────────────────────
    // CopySync 的 A 端必须拿到源文件路径才能在对端粘贴时回源。
    // 这一步几乎必然触发提示——关键是用户授权后能否永久静默。
    BOOL hasFile = [types containsObject:NSPasteboardTypeFileURL];
    if (hasFile) {
        t = beginStage(4, "读 public.file-url 实际内容 ← A 端必须拿到的东西");
        NSArray<NSURL *> *urls = [pb readObjectsForClasses:@[NSURL.class] options:nil];
        line("%lu 个文件：", (unsigned long)urls.count);
        for (NSURL *u in urls) {
            NSNumber *isDir = nil, *size = nil;
            [u getResourceValue:&isDir forKey:NSURLIsDirectoryKey error:nil];
            [u getResourceValue:&size  forKey:NSURLFileSizeKey    error:nil];
            fprintf(stderr, "             · %s%s  %lld bytes\n", u.path.UTF8String,
                    isDir.boolValue ? "/ (目录)" : "", size.longLongValue);
        }
        endStage(t);
    } else if ([types containsObject:NSPasteboardTypeString]) {
        t = beginStage(4, "读 public.utf8-plain-text 实际内容");
        NSString *s = [pb stringForType:NSPasteboardTypeString] ?: @"";
        NSString *p = s.length > 50 ? [[s substringToIndex:50] stringByAppendingString:@"…"] : s;
        line("%lu 字符：\"%s\"", (unsigned long)s.length,
             [p stringByReplacingOccurrencesOfString:@"\n" withString:@"⏎"].UTF8String);
        endStage(t);
    } else if ([types containsObject:NSPasteboardTypePNG] ||
               [types containsObject:NSPasteboardTypeTIFF]) {
        t = beginStage(4, "读图片位图数据");
        NSPasteboardType ty = [types containsObject:NSPasteboardTypePNG]
                            ? NSPasteboardTypePNG : NSPasteboardTypeTIFF;
        line("%s → %lu bytes", ty.UTF8String, (unsigned long)[pb dataForType:ty].length);
        endStage(t);
    }

    // ── 阶段 5：读取后再查授权状态 ─────────────────────────────
    // 文档称：首次触发提示后，behavior 会从 Default 自动变为 Ask，
    // 并开始出现在系统设置中。验证这个状态迁移。
    beginStage(5, "读取后再查 accessBehavior（观察状态迁移）");
    reportAccessBehavior(pb);

    fprintf(stderr, "\n  ✓ 本轮结束\n");
    fflush(stderr);
}

int main(int argc, const char *argv[]) {
    @autoreleasepool {
        NSApplication *app = NSApplication.sharedApplication;
        [app setActivationPolicy:NSApplicationActivationPolicyAccessory];

        // 轮数限制：便于脚本化测试，跑完 N 轮自动退出
        int maxRounds = (argc > 1) ? atoi(argv[1]) : 0;   // 0 = 无限

        banner("Spike A — 剪贴板隐私机制探测");
        NSPasteboard *pb = NSPasteboard.generalPasteboard;
        fprintf(stderr, "\n  bundle id : %s\n",
                NSBundle.mainBundle.bundleIdentifier.UTF8String ?: "(无)");
        BOOL previewOn = [NSUserDefaults.standardUserDefaults
                          boolForKey:@"EnablePasteboardPrivacyDeveloperPreview"];
        fprintf(stderr, "  隐私预览  : %s\n", previewOn ? "已开启 ✓（预演未来 macOS 行为）" : "未开启");
        fprintf(stderr, "  轮数限制  : %d\n", maxRounds);
        fflush(stderr);
        reportAccessBehavior(pb);

        // 轮询必须跑在主线程。
        //
        // 隐私机制开启后，剪贴板 API 会经由 ViewBridge 与系统 UI 服务通信
        // （为了显示授权弹窗），其内部的 NSCFRunLoopSemaphore 依赖 run loop，
        // 在裸 dispatch queue 上调用会抛异常并崩溃。实测堆栈：
        //   reportAccessBehavior → AppKit → ViewBridge -[NSCFRunLoopSemaphore wait]
        //                        → NSException raise → SIGSEGV
        //
        // 这条约束对 CopySync 同样成立：Go daemon 必须把剪贴板 cgo 调用
        // 调度到主线程，不能在任意 goroutine 里直接调。
        __block NSInteger last = pb.changeCount;
        __block int rounds = 0;
        __block BOOL busy = NO;   // probe 可能阻塞在授权弹窗上，防止定时器重入

        fprintf(stderr, "\n  已就绪（changeCount=%ld，主线程轮询），等待复制动作…\n", (long)last);
        fflush(stderr);

        [NSTimer scheduledTimerWithTimeInterval:kPollInterval repeats:YES block:^(NSTimer *timer) {
            if (busy) return;
            NSInteger cur = pb.changeCount;
            if (cur == last) return;
            last = cur;

            busy = YES;
            probe(pb, cur);
            busy = NO;

            if (maxRounds > 0 && ++rounds >= maxRounds) {
                fprintf(stderr, "\n  已完成 %d 轮，退出。\n", rounds);
                fflush(stderr);
                exit(0);
            }
        }];

        [app run];
    }
    return 0;
}
