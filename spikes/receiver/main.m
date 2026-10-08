// Spike D — promise 消费方（模拟「按下 Cmd+V」）
//
// NSFilePromiseReceiver 正是 Finder 等 App 在粘贴 file promise 时内部使用的类。
// 用它可以精确模拟粘贴行为，从而**完全自动化**地验证：
//
//   1. promise 回调是否真的被触发？
//   2. 回调里阻塞 N 秒是否会被系统超时杀掉？  ← 决定大文件传输是否可行
//   3. 文件内容是否完整送达？
//
// 这样不必依赖人工去各个 App 里按 Cmd+V，也不需要辅助功能权限。
//
// 用法：SpikeD [目标目录]

#import <Cocoa/Cocoa.h>

int main(int argc, const char *argv[]) {
    @autoreleasepool {
        NSApplication *app = NSApplication.sharedApplication;
        [app setActivationPolicy:NSApplicationActivationPolicyAccessory];

        NSString *destPath = argc > 1 ? @(argv[1]) : @"/tmp/copysync-paste-dest";
        NSURL *dest = [NSURL fileURLWithPath:destPath];
        [NSFileManager.defaultManager createDirectoryAtURL:dest
                               withIntermediateDirectories:YES attributes:nil error:nil];

        NSPasteboard *pb = NSPasteboard.generalPasteboard;

        fprintf(stderr, "\n[SpikeD] bundle=%s\n",
                NSBundle.mainBundle.bundleIdentifier.UTF8String);
        fprintf(stderr, "[SpikeD] 粘贴目标目录：%s\n", destPath.UTF8String);
        fprintf(stderr, "[SpikeD] 剪贴板类型：\n");
        for (NSPasteboardType t in pb.types) {
            fprintf(stderr, "           · %s\n", t.UTF8String);
        }

        // 从剪贴板取出 promise。这一步等价于 Finder 识别出「剪贴板里是承诺文件」。
        NSArray *receivers = [pb readObjectsForClasses:@[NSFilePromiseReceiver.class]
                                               options:nil];
        if (receivers.count == 0) {
            fprintf(stderr, "[SpikeD] ✗ 剪贴板里没有 file promise\n");
            fprintf(stderr, "         （若剪贴板是普通文件 URL，说明写入方没用 promise 机制）\n");
            return 1;
        }

        fprintf(stderr, "[SpikeD] ✓ 找到 %lu 个 promise，开始请求文件…\n\n",
                (unsigned long)receivers.count);

        NSOperationQueue *q = NSOperationQueue.new;
        __block NSInteger pending = receivers.count;
        __block NSInteger failed = 0;
        NSDate *t0 = NSDate.date;

        for (NSFilePromiseReceiver *r in receivers) {
            fprintf(stderr, "[SpikeD] promise 声明的文件名：%s\n",
                    r.fileNames.count ? r.fileNames.description.UTF8String : "(未提供)");

            // ★ 这一步触发对端的 writePromiseToURL 回调，等价于用户按下 Cmd+V
            [r receivePromisedFilesAtDestination:dest
                                         options:@{}
                                  operationQueue:q
                                          reader:^(NSURL *fileURL, NSError *err) {
                NSTimeInterval dt = -[t0 timeIntervalSinceNow];
                if (err) {
                    failed++;
                    fprintf(stderr, "[SpikeD] ✗ 失败（%.1fs）：%s\n",
                            dt, err.localizedDescription.UTF8String);
                } else {
                    NSNumber *size = nil;
                    [fileURL getResourceValue:&size forKey:NSURLFileSizeKey error:nil];
                    NSNumber *isDir = nil;
                    [fileURL getResourceValue:&isDir forKey:NSURLIsDirectoryKey error:nil];
                    fprintf(stderr, "[SpikeD] ✓ 收到（耗时 %.1fs）：%s  %lld bytes%s\n",
                            dt, fileURL.path.UTF8String, size.longLongValue,
                            isDir.boolValue ? "  (目录)" : "");
                    // 打印内容前 200 字节，确认数据完整
                    if (!isDir.boolValue) {
                        NSString *body = [NSString stringWithContentsOfURL:fileURL
                                                                  encoding:NSUTF8StringEncoding
                                                                     error:nil];
                        if (body.length) {
                            NSString *p = body.length > 200 ? [body substringToIndex:200] : body;
                            fprintf(stderr, "           内容：%s\n",
                                    [p stringByReplacingOccurrencesOfString:@"\n"
                                                                 withString:@" ⏎ "].UTF8String);
                        }
                    }
                }
                if (--pending == 0) {
                    fprintf(stderr, "\n[SpikeD] 全部完成，总耗时 %.1fs，失败 %ld 个\n",
                            -[t0 timeIntervalSinceNow], (long)failed);
                    fflush(stderr);
                    exit(failed ? 1 : 0);
                }
            }];
        }

        // 兜底超时：promise 若被系统静默丢弃，这里会兜住
        dispatch_after(dispatch_time(DISPATCH_TIME_NOW, (int64_t)(300 * NSEC_PER_SEC)),
                       dispatch_get_main_queue(), ^{
            fprintf(stderr, "\n[SpikeD] ✗ 300 秒仍未全部完成，判定超时\n");
            exit(2);
        });

        [app run];
    }
    return 0;
}
