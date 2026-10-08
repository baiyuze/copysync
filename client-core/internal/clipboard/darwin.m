//go:build darwin

// macOS 剪贴板的原生实现。
//
// ⚠️ 本文件所有函数都必须在主线程调用（见 clipboard.go 的线程约束说明）。
// 隐私机制开启后，剪贴板 API 经 ViewBridge 与系统 UI 服务通信，
// 其内部 NSCFRunLoopSemaphore 依赖 run loop，在裸线程上调用会直接 SIGSEGV。

#import <Cocoa/Cocoa.h>
#import <UniformTypeIdentifiers/UniformTypeIdentifiers.h>

#include <stdlib.h>
#include <string.h>

// 返回给 Go 的字符串统一用 strdup 分配，由 Go 侧 C.free 释放。
static char *copyCString(NSString *s) {
    if (s == nil) return NULL;
    const char *utf8 = s.UTF8String;
    if (utf8 == NULL) return NULL;
    return strdup(utf8);
}

static char *copyJSON(id obj) {
    if (obj == nil) return NULL;
    NSError *err = nil;
    NSData *data = [NSJSONSerialization dataWithJSONObject:obj options:0 error:&err];
    if (err != nil || data == nil) return NULL;
    return copyCString([[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding]);
}

// 在主线程等待异步回调：spin run loop 而不是阻塞。
//
// 直接阻塞会有两个问题：completionHandler 若投递到主队列则死锁；
// 且系统授权弹窗依赖主 run loop 才能显示和响应。
static BOOL waitSpinning(dispatch_semaphore_t sem, NSTimeInterval timeout) {
    NSDate *deadline = [NSDate dateWithTimeIntervalSinceNow:timeout];
    while (dispatch_semaphore_wait(sem, DISPATCH_TIME_NOW) != 0) {
        if (deadline.timeIntervalSinceNow < 0) return NO;
        [NSRunLoop.currentRunLoop runMode:NSDefaultRunLoopMode
                               beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.01]];
    }
    return YES;
}

#pragma mark - 免提示探测

long cs_change_count(void) {
    return (long)NSPasteboard.generalPasteboard.changeCount;
}

// 返回值对应 Go 侧的 Permission 常量：
// -1 = 系统无此机制，0 = default，1 = ask，2 = alwaysAllow，3 = alwaysDeny
int cs_access_behavior(void) {
    if (@available(macOS 15.4, *)) {
        return (int)NSPasteboard.generalPasteboard.accessBehavior;
    }
    return -1;
}

// 类型列表。实测免授权，隐私机制开启时耗时约 0.003s。
char *cs_types_json(void) {
    NSArray<NSPasteboardType> *types = NSPasteboard.generalPasteboard.types;
    return copyJSON(types ?: @[]);
}

// 用 detectMetadataForTypes 探测内容类型。
//
// Apple 明确保证此方法不会通知用户，实测隐私机制开启时耗时 0.021s 且无弹窗。
// CopySync 靠它在不触发授权的前提下完成"复制的是文件还是文本"的分类。
char *cs_detect_content_type(void) {
    if (@available(macOS 15.4, *)) {
        __block char *result = NULL;
        dispatch_semaphore_t sem = dispatch_semaphore_create(0);

        [NSPasteboard.generalPasteboard
            detectMetadataForTypes:[NSSet setWithObject:NSPasteboardMetadataTypeContentType]
                 completionHandler:^(NSDictionary<NSPasteboardMetadataType, id> *md, NSError *err) {
            if (err == nil) {
                UTType *ct = md[NSPasteboardMetadataTypeContentType];
                if (ct != nil) result = copyCString(ct.identifier);
            }
            dispatch_semaphore_signal(sem);
        }];

        waitSpinning(sem, 5.0);
        return result;
    }
    return NULL;
}

#pragma mark - 读取内容（可能触发授权）

char *cs_read_text(void) {
    NSPasteboard *pb = NSPasteboard.generalPasteboard;

    // 依次尝试多种取法。
    //
    // stringForType: 对某些写入方（例如只声明了旧式 NSStringPboardType 的应用）
    // 会返回 nil，而 readObjectsForClasses: 能正确走类型转换拿到内容。
    // 剪贴板的写入方五花八门，这里不假定它们都规范声明了现代 UTI。
    NSString *s = [pb stringForType:NSPasteboardTypeString];
    if (s.length > 0) return copyCString(s);

    NSArray<NSString *> *objs =
        [pb readObjectsForClasses:@[NSString.class] options:nil];
    if (objs.count > 0) {
        NSString *first = objs.firstObject;
        if (first.length > 0) return copyCString(first);
    }

    // 最后退回原始字节：个别写入方提供的数据不带 BOM/编码声明，
    // 高层 API 会拒绝转换，但按 UTF-8 解出来往往是对的。
    NSData *data = [pb dataForType:NSPasteboardTypeString];
    if (data.length > 0) {
        NSString *decoded = [[NSString alloc] initWithData:data
                                                  encoding:NSUTF8StringEncoding];
        if (decoded.length > 0) return copyCString(decoded);
    }
    return NULL;
}

char *cs_read_html(void) {
    NSPasteboard *pb = NSPasteboard.generalPasteboard;
    NSString *s = [pb stringForType:NSPasteboardTypeHTML];
    if (s.length > 0) return copyCString(s);

    NSData *data = [pb dataForType:NSPasteboardTypeHTML];
    if (data.length > 0) {
        NSString *decoded = [[NSString alloc] initWithData:data
                                                  encoding:NSUTF8StringEncoding];
        if (decoded.length > 0) return copyCString(decoded);
    }
    return NULL;
}

// 返回文件路径数组的 JSON。CopySync 的 A 端靠它拿到源文件位置。
char *cs_read_file_urls_json(void) {
    NSArray<NSURL *> *urls = [NSPasteboard.generalPasteboard
        readObjectsForClasses:@[NSURL.class]
                      options:@{NSPasteboardURLReadingFileURLsOnlyKey: @YES}];
    if (urls == nil) return NULL;

    NSMutableArray<NSString *> *paths = [NSMutableArray arrayWithCapacity:urls.count];
    for (NSURL *u in urls) {
        if (u.isFileURL && u.path != nil) [paths addObject:u.path];
    }
    return copyJSON(paths);
}

// 读取图片并统一编码为 PNG。
// 剪贴板里可能是 TIFF（系统截图常见），统一成 PNG 便于跨平台传输。
int cs_read_image_png(void **out, int *outLen) {
    *out = NULL;
    *outLen = 0;

    NSPasteboard *pb = NSPasteboard.generalPasteboard;
    NSData *png = [pb dataForType:NSPasteboardTypePNG];

    if (png == nil) {
        NSData *tiff = [pb dataForType:NSPasteboardTypeTIFF];
        if (tiff == nil) return 0;
        NSBitmapImageRep *rep = [[NSBitmapImageRep alloc] initWithData:tiff];
        if (rep == nil) return 0;
        png = [rep representationUsingType:NSBitmapImageFileTypePNG properties:@{}];
        if (png == nil) return 0;
    }

    void *buf = malloc(png.length);
    if (buf == NULL) return 0;
    memcpy(buf, png.bytes, png.length);
    *out = buf;
    *outLen = (int)png.length;
    return 1;
}

#pragma mark - 写入剪贴板

int cs_write_text(const char *utf8) {
    if (utf8 == NULL) return 0;
    NSString *s = @(utf8);
    NSPasteboard *pb = NSPasteboard.generalPasteboard;
    [pb clearContents];
    return [pb setString:s forType:NSPasteboardTypeString] ? 1 : 0;
}

int cs_write_html(const char *html, const char *plain) {
    if (html == NULL) return 0;
    NSPasteboard *pb = NSPasteboard.generalPasteboard;
    [pb clearContents];
    BOOL ok = [pb setString:@(html) forType:NSPasteboardTypeHTML];
    // 同时放纯文本兜底：很多应用不接受 HTML
    if (plain != NULL) [pb setString:@(plain) forType:NSPasteboardTypeString];
    return ok ? 1 : 0;
}

int cs_write_image_png(const void *data, int len) {
    if (data == NULL || len <= 0) return 0;
    NSData *png = [NSData dataWithBytes:data length:(NSUInteger)len];
    NSPasteboard *pb = NSPasteboard.generalPasteboard;
    [pb clearContents];
    return [pb setData:png forType:NSPasteboardTypePNG] ? 1 : 0;
}

// 把本地文件写进剪贴板，使其可被 Finder 直接粘贴。
//
// 文件必须已真实存在于磁盘——M-1 验证过，延迟渲染（promise）在 macOS 的
// 剪贴板粘贴路径上不成立，所以这里只能写真实路径。
int cs_write_files_json(const char *pathsJSON) {
    if (pathsJSON == NULL) return 0;

    NSData *data = [@(pathsJSON) dataUsingEncoding:NSUTF8StringEncoding];
    NSArray *paths = [NSJSONSerialization JSONObjectWithData:data options:0 error:nil];
    if (![paths isKindOfClass:NSArray.class] || paths.count == 0) return 0;

    NSMutableArray<NSURL *> *urls = [NSMutableArray arrayWithCapacity:paths.count];
    for (NSString *p in paths) {
        if (![p isKindOfClass:NSString.class]) continue;
        // 不存在的路径写进去只会让粘贴失败，不如提前剔除
        if (![NSFileManager.defaultManager fileExistsAtPath:p]) continue;
        [urls addObject:[NSURL fileURLWithPath:p]];
    }
    if (urls.count == 0) return 0;

    NSPasteboard *pb = NSPasteboard.generalPasteboard;
    [pb clearContents];
    return [pb writeObjects:urls] ? 1 : 0;
}

#pragma mark - 权限引导

void cs_open_permission_settings(void) {
    // 直达「隐私与安全性 → 从其他 App 粘贴」
    NSURL *url = [NSURL URLWithString:
        @"x-apple.systempreferences:com.apple.preference.security?Privacy_Pasteboard"];
    [NSWorkspace.sharedWorkspace openURL:url];
}

#pragma mark - 主线程 run loop

// 驱动 run loop 一小段时间。
//
// daemon 的主线程需要持续转动 run loop，否则剪贴板 API 的跨进程回调
// 与授权弹窗都无法工作。Go 侧在锁定的主线程上循环调用它。
void cs_run_loop_once(double seconds) {
    @autoreleasepool {
        [NSRunLoop.currentRunLoop runMode:NSDefaultRunLoopMode
                               beforeDate:[NSDate dateWithTimeIntervalSinceNow:seconds]];
    }
}

// 初始化 NSApplication。
//
// 必须有 NSApplication 实例，AppKit 的剪贴板与 ViewBridge 机制才完整可用；
// Accessory 模式让 daemon 不出现在 Dock 和任务切换器里。
void cs_init_app(void) {
    @autoreleasepool {
        NSApplication *app = NSApplication.sharedApplication;
        [app setActivationPolicy:NSApplicationActivationPolicyAccessory];
    }
}
