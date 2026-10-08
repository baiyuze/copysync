// Spike C — 「复制方」模拟器
//
// 存在的唯一理由：给 Spike A 提供一个**真实 GUI App** 身份的剪贴板写入方。
//
// pbcopy / osascript 是命令行工具，TCC 可能不把它们当作「其他 App」，
// 用它们测出的「没有隐私提示」不可采信。这个程序有自己独立的
// bundle identifier 且是 NSApplication，与真实 App（Finder/微信）同等地位。
//
// 用法：
//   SpikeC file <路径>    把文件 URL 写进剪贴板
//   SpikeC text <内容>    把文本写进剪贴板
//   SpikeC image          把一张生成的图片写进剪贴板

#import <Cocoa/Cocoa.h>

int main(int argc, const char *argv[]) {
    @autoreleasepool {
        NSApplication *app = NSApplication.sharedApplication;
        [app setActivationPolicy:NSApplicationActivationPolicyAccessory];

        NSString *mode = argc > 1 ? @(argv[1]) : @"text";
        NSString *arg  = argc > 2 ? @(argv[2]) : @"default payload";

        NSPasteboard *pb = NSPasteboard.generalPasteboard;
        [pb clearContents];

        BOOL ok = NO;
        NSString *desc = nil;

        if ([mode isEqualToString:@"file"]) {
            NSURL *url = [NSURL fileURLWithPath:arg];
            ok = [pb writeObjects:@[url]];
            desc = [NSString stringWithFormat:@"file-url → %@", arg];

        } else if ([mode isEqualToString:@"image"]) {
            // 生成一张纯色图，模拟截图
            NSImage *img = [[NSImage alloc] initWithSize:NSMakeSize(320, 200)];
            [img lockFocus];
            [[NSColor systemTealColor] setFill];
            NSRectFill(NSMakeRect(0, 0, 320, 200));
            [img unlockFocus];
            ok = [pb writeObjects:@[img]];
            desc = @"image (320x200)";

        } else {
            ok = [pb setString:arg forType:NSPasteboardTypeString];
            desc = [NSString stringWithFormat:@"text → \"%@\"", arg];
        }

        fprintf(stderr,
                "[SpikeC] bundle=%s\n"
                "[SpikeC] 已写入：%s\n"
                "[SpikeC] 结果=%s  changeCount=%ld\n"
                "[SpikeC] 保持存活 40 秒（模拟真实 App 持有剪贴板）…\n",
                NSBundle.mainBundle.bundleIdentifier.UTF8String,
                desc.UTF8String,
                ok ? "OK" : "FAIL",
                (long)pb.changeCount);
        fflush(stderr);

        // 真实 App 复制后不会立刻退出。保持存活，确保读方读取时
        // 剪贴板 owner 仍是一个活着的、不同 bundle id 的 GUI 进程。
        [NSThread sleepForTimeInterval:40];
    }
    return 0;
}
