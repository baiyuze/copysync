// Spike B — NSFilePromiseProvider 可用性与兼容性验证
//
// 目的：验证 CopySync 的核心机制是否成立——
//   1. 往剪贴板写 file promise 后，目标 App 粘贴时能否触发回调？
//      （这是「感知到对方按下 Cmd+V」的唯一途径）
//   2. 回调里能否长时间阻塞（模拟跨网络拉取）而不被系统超时杀掉？
//      （决定大文件传输是否可行）
//   3. 哪些 App 支持 promise、哪些只认传统 file-url？
//      （决定是否需要降级策略）
//
// 用法：运行后按提示输入命令写入不同形态的剪贴板内容，
//      然后去各个 App 里 Cmd+V，观察终端日志。

#import <Cocoa/Cocoa.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>

static NSMutableArray *gKeepAlive;          // 保持 provider 存活，否则 ARC 会提前释放
static NSOperationQueue *gWorkQueue;        // promise 回调在此队列执行，可安全阻塞
static NSTimeInterval gFakeLatency = 5.0;   // 模拟「跨网络拉取」的耗时
static NSDate *gWriteTime;                  // 写入剪贴板的时刻，用于计算用户多久后粘贴

static void banner(const char *s) {
    fprintf(stderr, "\n════════════════════════════════════════════════════════\n");
    fprintf(stderr, " %s\n", s);
    fprintf(stderr, "════════════════════════════════════════════════════════\n");
    fflush(stderr);
}

static void say(const char *fmt, ...) {
    va_list ap; va_start(ap, fmt);
    vfprintf(stderr, fmt, ap);
    fprintf(stderr, "\n");
    va_end(ap);
    fflush(stderr);
}

#pragma mark - Promise delegate

@interface SpikeDelegate : NSObject <NSFilePromiseProviderDelegate>
@end

@implementation SpikeDelegate

// 系统询问「这个 promise 对应的文件叫什么名字」。
// 注意：这一步在粘贴发生之前就会被调用（目标 App 要先知道文件名），
// 所以它**不能**作为「用户按下了 Cmd+V」的信号。
- (NSString *)filePromiseProvider:(NSFilePromiseProvider *)provider
                   fileNameForType:(NSString *)fileType {
    NSString *name = provider.userInfo[@"name"];
    say("  ├─ [fileNameForType] 被询问文件名 → %s   (type=%s)",
        name.UTF8String, fileType.UTF8String);
    say("     ⚠️  注意：此回调在粘贴前就会触发，不能当作粘贴信号");
    return name;
}

// ★ 这才是真正的「用户按下了 Cmd+V」信号。
// 系统给我们一个目标 URL，我们在这里把数据写进去。
// CopySync 正式实现中，这里会发出 PullRequest 并把网络流直接写入 url。
- (void)filePromiseProvider:(NSFilePromiseProvider *)provider
          writePromiseToURL:(NSURL *)url
          completionHandler:(void (^)(NSError *))done {
    NSString *name = provider.userInfo[@"name"];
    BOOL isFolder = [provider.userInfo[@"folder"] boolValue];
    NSTimeInterval sincePaste = -[gWriteTime timeIntervalSinceNow];

    say("\n  ★★★ [writePromiseToURL] 粘贴被触发！");
    say("  ├─ 文件      : %s%s", name.UTF8String, isFolder ? " (目录)" : "");
    say("  ├─ 目标路径  : %s", url.path.UTF8String);
    say("  ├─ 距写入剪贴板 %.1f 秒", sincePaste);
    say("  ├─ 所在线程  : %s", NSThread.isMainThread ? "主线程 ⚠️" : "后台线程 ✓");
    say("  ├─ 模拟拉取中，将阻塞 %.0f 秒……", gFakeLatency);

    // 模拟跨网络拉取。真实实现中这里是 DataChannel → zstd → tar 的流式写入。
    NSDate *t0 = NSDate.date;
    fprintf(stderr, "  │  ");
    for (NSTimeInterval spent = 0; spent < gFakeLatency; spent += 0.5) {
        [NSThread sleepForTimeInterval:0.5];
        fprintf(stderr, "▓"); fflush(stderr);
    }
    fprintf(stderr, "\n");

    NSError *err = nil;
    NSString *body = [NSString stringWithFormat:
        @"CopySync spike B\n写入时刻: %@\n阻塞时长: %.1f 秒\n用户粘贴延迟: %.1f 秒\n",
        NSDate.date, gFakeLatency, sincePaste];

    if (isFolder) {
        [NSFileManager.defaultManager createDirectoryAtURL:url
                               withIntermediateDirectories:YES
                                                attributes:nil error:&err];
        if (!err) {
            for (int i = 1; i <= 3; i++) {
                NSURL *f = [url URLByAppendingPathComponent:
                            [NSString stringWithFormat:@"file-%d.txt", i]];
                [body writeToURL:f atomically:YES encoding:NSUTF8StringEncoding error:&err];
            }
        }
    } else {
        [body writeToURL:url atomically:YES encoding:NSUTF8StringEncoding error:&err];
    }

    say("  ├─ 写入完成，实际耗时 %.1f 秒", -[t0 timeIntervalSinceNow]);
    say("  └─ %s", err ? (const char *)[NSString stringWithFormat:@"失败: %@", err].UTF8String
                       : "成功 ✓ —— 去目标 App 确认文件是否真的出现了");
    done(err);
}

// 指定回调执行的队列。必须是后台队列，否则阻塞会冻结整个 App。
- (NSOperationQueue *)operationQueueForFilePromiseProvider:(NSFilePromiseProvider *)p {
    return gWorkQueue;
}

@end

#pragma mark - 延迟提供 file-url（老式 pasteboard owner 机制）

// NSFilePromiseProvider 主要服务于拖放；Finder 的剪贴板粘贴只认
// public.file-url / NSFilenamesPboardType。
//
// 老式的 declareTypes:owner: 机制允许我们声明这些 Finder 认识的类型，
// 但**不立即提供数据**——系统会在有人真正来取时回调 provideDataForType:。
// 这同样能拿到「用户按下了 Cmd+V」的信号，且目标 App 的兼容性好得多。
@interface LazyFileOwner : NSObject
@end

@implementation LazyFileOwner

- (void)pasteboard:(NSPasteboard *)sender provideDataForType:(NSPasteboardType)type {
    NSTimeInterval since = -[gWriteTime timeIntervalSinceNow];
    say("\n  ★★★ [provideDataForType] 粘贴被触发！");
    say("  ├─ 请求类型  : %s", type.UTF8String);
    say("  ├─ 距写入剪贴板 %.1f 秒", since);
    say("  ├─ 所在线程  : %s", NSThread.isMainThread ? "主线程 ← 注意会冻结 UI" : "后台线程");
    say("  ├─ 模拟跨网络拉取，阻塞 %.0f 秒……", gFakeLatency);

    NSDate *t0 = NSDate.date;
    fprintf(stderr, "  │  ");
    for (NSTimeInterval s = 0; s < gFakeLatency; s += 0.5) {
        [NSThread sleepForTimeInterval:0.5];
        fprintf(stderr, "▓"); fflush(stderr);
    }
    fprintf(stderr, "\n");

    // 此刻才落盘。真实实现中这里是 DataChannel → zstd → tar 的流式写入。
    NSString *dir = [NSTemporaryDirectory() stringByAppendingPathComponent:@"copysync-lazy"];
    [NSFileManager.defaultManager createDirectoryAtPath:dir
                            withIntermediateDirectories:YES attributes:nil error:nil];
    NSString *path = [dir stringByAppendingPathComponent:@"copysync-lazy.txt"];
    NSError *err = nil;
    [[NSString stringWithFormat:
        @"CopySync — 延迟粘贴验证\n落盘时刻: %@\n阻塞时长: %.1f 秒\n粘贴延迟: %.1f 秒\n",
        NSDate.date, gFakeLatency, since]
     writeToFile:path atomically:YES encoding:NSUTF8StringEncoding error:&err];

    BOOL ok = NO;
    if ([type isEqualToString:NSPasteboardTypeFileURL]) {
        ok = [sender setString:[NSURL fileURLWithPath:path].absoluteString forType:type];
    } else if ([type isEqualToString:NSFilenamesPboardType]) {
        ok = [sender setPropertyList:@[path] forType:type];
    } else {
        ok = [sender setString:path forType:type];
    }

    say("  ├─ 落盘%s  提供数据%s", err ? "失败" : "成功", ok ? "成功" : "失败");
    say("  └─ 回调耗时 %.1f 秒 —— 去 Finder 确认文件是否出现", -[t0 timeIntervalSinceNow]);
}

@end

#pragma mark - 写入剪贴板的几种形态

static SpikeDelegate *gDelegate;
static LazyFileOwner *gLazyOwner;

// 声明 Finder 认识的类型，但不提供数据——等对方来取时才回调
static void writeLazyFileURL(void) {
    NSPasteboard *pb = NSPasteboard.generalPasteboard;
    // declareTypes:owner: 自身会清空剪贴板并递增 changeCount
    [pb declareTypes:@[NSPasteboardTypeFileURL, NSFilenamesPboardType] owner:gLazyOwner];
    gWriteTime = NSDate.date;

    say("\n  ✓ 已写入剪贴板：延迟 file-url（declareTypes:owner:，数据暂不提供）");
    say("    剪贴板类型：");
    for (NSPasteboardType t in pb.types) fprintf(stderr, "      · %s\n", t.UTF8String);
    say("    现在去 Finder 按 Cmd+V");
}

static NSFilePromiseProvider *makeProvider(NSString *name, BOOL folder) {
    UTType *type = folder ? UTTypeFolder
                          : ([UTType typeWithFilenameExtension:name.pathExtension] ?: UTTypeData);
    NSFilePromiseProvider *p = [[NSFilePromiseProvider alloc] initWithFileType:type.identifier
                                                                      delegate:gDelegate];
    p.userInfo = @{@"name": name, @"folder": @(folder)};
    [gKeepAlive addObject:p];
    return p;
}

static void writeToPasteboard(NSArray *objects, const char *what) {
    NSPasteboard *pb = NSPasteboard.generalPasteboard;
    [pb clearContents];
    BOOL ok = [pb writeObjects:objects];
    gWriteTime = NSDate.date;

    say("\n  ✓ 已写入剪贴板：%s", what);
    say("    writeObjects 返回 %s", ok ? "YES" : "NO");
    say("    剪贴板类型：");
    for (NSPasteboardType t in pb.types) {
        fprintf(stderr, "      · %s\n", t.UTF8String);
    }
    say("    现在去目标 App 按 Cmd+V，观察是否触发 writePromiseToURL");
}

// 对照组：传统的真实文件 URL。用来区分「App 不支持 promise」和「App 完全不能粘文件」。
static void writeRealFile(void) {
    NSURL *dir = [NSURL fileURLWithPath:NSTemporaryDirectory()];
    NSURL *f = [dir URLByAppendingPathComponent:@"copysync-real-file.txt"];
    [@"这是一个真实存在的文件（对照组）\n" writeToURL:f atomically:YES
                                        encoding:NSUTF8StringEncoding error:nil];
    writeToPasteboard(@[f], "对照组 — 传统 file-url（文件已真实存在）");
    say("    路径：%s", f.path.UTF8String);
}

#pragma mark - main

// 执行一条写入命令。交互模式和命令行模式共用。
// 返回 NO 表示命令无法识别。
static BOOL doCommand(NSString *cmd) {
    if ([cmd isEqualToString:@"1"] || [cmd isEqualToString:@"single"]) {
        writeToPasteboard(@[makeProvider(@"copysync-single.txt", NO)], "单个文件 promise");
    } else if ([cmd isEqualToString:@"2"] || [cmd isEqualToString:@"folder"]) {
        writeToPasteboard(@[makeProvider(@"copysync-folder", YES)], "文件夹 promise");
    } else if ([cmd isEqualToString:@"3"] || [cmd isEqualToString:@"multi"]) {
        writeToPasteboard(@[makeProvider(@"copysync-a.txt", NO),
                            makeProvider(@"copysync-b.txt", NO),
                            makeProvider(@"copysync-c.txt", NO)],
                          "多文件 promise（3 个）");
    } else if ([cmd isEqualToString:@"4"] || [cmd isEqualToString:@"real"]) {
        writeRealFile();
    } else if ([cmd isEqualToString:@"5"] || [cmd isEqualToString:@"lazy"]) {
        writeLazyFileURL();
    } else {
        return NO;
    }
    return YES;
}

static void printMenu(void) {
    say("\n  ┌── 命令 ─────────────────────────────────────────");
    say("  │  1      单个文件 promise (.txt)");
    say("  │  2      文件夹 promise（内含 3 个文件）");
    say("  │  3      多文件 promise（3 个 .txt）");
    say("  │  4      对照组：传统 file-url（真实文件）");
    say("  │  5      延迟 file-url（declareTypes:owner:）← 最有希望的方案");
    say("  │  d<秒>  设置模拟拉取延迟，如 d5 / d30 / d120");
    say("  │  q      退出");
    say("  └─────────────────────────────────────────────────");
    fprintf(stderr, "  当前延迟 %.0f 秒 > ", gFakeLatency);
    fflush(stderr);
}

int main(int argc, const char *argv[]) {
    @autoreleasepool {
        gKeepAlive = NSMutableArray.new;
        gDelegate = SpikeDelegate.new;
        gLazyOwner = LazyFileOwner.new;   // 必须保持存活，否则系统回调时 owner 已释放
        gWorkQueue = NSOperationQueue.new;
        gWorkQueue.name = @"copysync.spike.promise";
        gWorkQueue.maxConcurrentOperationCount = 4;

        NSApplication *app = NSApplication.sharedApplication;
        [app setActivationPolicy:NSApplicationActivationPolicyAccessory];

        banner("Spike B — NSFilePromiseProvider 可用性验证");

        // 命令行模式：SpikeB <mode> [延迟秒数]
        // 写入后保持存活等待消费方，供脚本自动化测试。
        if (argc > 1) {
            if (argc > 2) gFakeLatency = atof(argv[2]);
            NSString *mode = @(argv[1]);
            say("  模式：%s   模拟拉取延迟：%.0f 秒", mode.UTF8String, gFakeLatency);
            if (!doCommand(mode)) {
                say("  ✗ 未知模式：%s（可用：single/folder/multi/real）", mode.UTF8String);
                return 1;
            }
            say("\n  等待消费方取件…（进程保持存活）");
            [app run];
            return 0;
        }

        say("\n  要回答的问题：");
        say("    1. 粘贴时能否触发 writePromiseToURL 回调？");
        say("    2. 回调里阻塞 %.0f 秒会被系统超时杀掉吗？", gFakeLatency);
        say("    3. 哪些 App 支持 promise？（Finder / 微信 / 飞书 / Mail / VS Code …）");
        say("\n  测试方法：选一个命令写入剪贴板 → 去目标 App 按 Cmd+V → 看日志");

        dispatch_async(dispatch_get_global_queue(QOS_CLASS_UTILITY, 0), ^{
            char buf[64];
            while (true) {
                printMenu();
                if (!fgets(buf, sizeof buf, stdin)) break;
                NSString *cmd = [[NSString stringWithUTF8String:buf]
                    stringByTrimmingCharactersInSet:NSCharacterSet.whitespaceAndNewlineCharacterSet];

                if ([cmd isEqualToString:@"q"]) { exit(0); }
                else if (doCommand(cmd)) { /* 已处理 */ }
                else if ([cmd hasPrefix:@"d"]) {
                    double v = [cmd substringFromIndex:1].doubleValue;
                    if (v > 0) { gFakeLatency = v; say("  ✓ 延迟设为 %.0f 秒", v); }
                    else say("  ✗ 无效数值");
                }
                else if (cmd.length) say("  ✗ 未知命令：%s", cmd.UTF8String);
            }
        });

        [app run];
    }
    return 0;
}
