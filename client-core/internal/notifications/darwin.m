//go:build darwin

#import <Cocoa/Cocoa.h>
#import <UserNotifications/UserNotifications.h>
#include <stdlib.h>
#include <string.h>
extern void goNotificationAction(char *, char *);

static BOOL spin(dispatch_semaphore_t sem) {
    NSDate *deadline = [NSDate dateWithTimeIntervalSinceNow:8];
    while (dispatch_semaphore_wait(sem, DISPATCH_TIME_NOW) != 0) {
        if (deadline.timeIntervalSinceNow < 0) return NO;
        @autoreleasepool {
            [NSRunLoop.currentRunLoop runMode:NSDefaultRunLoopMode beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.01]];
        }
    }
    return YES;
}

@interface PFNotificationDelegate : NSObject <UNUserNotificationCenterDelegate>
@end
@implementation PFNotificationDelegate
- (void)userNotificationCenter:(UNUserNotificationCenter *)center willPresentNotification:(UNNotification *)notification withCompletionHandler:(void (^)(UNNotificationPresentationOptions))completionHandler {
    completionHandler(UNNotificationPresentationOptionBanner | UNNotificationPresentationOptionList);
}
- (void)userNotificationCenter:(UNUserNotificationCenter *)center didReceiveNotificationResponse:(UNNotificationResponse *)response withCompletionHandler:(void (^)(void))completionHandler {
    NSString *verb = response.actionIdentifier;
    if ([verb isEqualToString:UNNotificationDefaultActionIdentifier]) verb = @"open";
    if (![verb isEqualToString:UNNotificationDismissActionIdentifier]) {
        goNotificationAction((char *)response.notification.request.identifier.UTF8String, (char *)verb.UTF8String);
    }
    completionHandler();
}
@end

static PFNotificationDelegate *delegate;
int pf_notify_init(void) {
    @autoreleasepool {
        if (NSBundle.mainBundle.bundleIdentifier.length == 0) return 0;
        delegate = [PFNotificationDelegate new];
        UNUserNotificationCenter.currentNotificationCenter.delegate = delegate;
        // Delivered notifications from an earlier daemon instance cannot retain
        // the state of a transfer. History remains the durable source of truth.
        [UNUserNotificationCenter.currentNotificationCenter removeAllDeliveredNotifications];
        return 1;
    }
}
int pf_notify_permission(int request) {
    @autoreleasepool {
        UNUserNotificationCenter *center = UNUserNotificationCenter.currentNotificationCenter;
        __block NSInteger status = -1;
        dispatch_semaphore_t sem = dispatch_semaphore_create(0);
        [center getNotificationSettingsWithCompletionHandler:^(UNNotificationSettings *settings) { status = settings.authorizationStatus; dispatch_semaphore_signal(sem); }];
        if (!spin(sem)) return -1;
        if (request && status == UNAuthorizationStatusNotDetermined) {
            sem = dispatch_semaphore_create(0);
            [center requestAuthorizationWithOptions:UNAuthorizationOptionAlert completionHandler:^(BOOL granted, NSError *error) { status = granted ? UNAuthorizationStatusAuthorized : UNAuthorizationStatusDenied; dispatch_semaphore_signal(sem); }];
            if (!spin(sem)) return -1;
        } else if (request && status != UNAuthorizationStatusNotDetermined) {
            [NSWorkspace.sharedWorkspace openURL:[NSURL URLWithString:@"x-apple.systempreferences:com.apple.Notifications-Settings.extension?id=com.pasteferry.daemon"]];
        }
        return (int)status;
    }
}
int pf_notify_show(const char *id, const char *title, const char *body, const char *action, const char *label) {
    @autoreleasepool {
        UNUserNotificationCenter *center = UNUserNotificationCenter.currentNotificationCenter;
        NSString *identifier = @(id), *verb = @(action);
        // Keep all categories registered: notifications for different records may
        // remain in Notification Center at the same time.
        UNNotificationAction *button = [UNNotificationAction actionWithIdentifier:verb title:@(label) options:UNNotificationActionOptionNone];
        UNNotificationCategory *category = [UNNotificationCategory categoryWithIdentifier:verb actions:@[button] intentIdentifiers:@[] options:UNNotificationCategoryOptionNone];
        dispatch_semaphore_t sem = dispatch_semaphore_create(0);
        [center getNotificationCategoriesWithCompletionHandler:^(NSSet<UNNotificationCategory *> *categories) {
            NSMutableSet *updated = [categories mutableCopy];
            for (UNNotificationCategory *old in categories) if ([old.identifier isEqualToString:verb]) [updated removeObject:old];
            [updated addObject:category]; [center setNotificationCategories:updated]; dispatch_semaphore_signal(sem);
        }];
        if (!spin(sem)) return 0;
        UNMutableNotificationContent *content = [UNMutableNotificationContent new];
        content.title = @(title); content.body = @(body); content.categoryIdentifier = verb;
        [center removeDeliveredNotificationsWithIdentifiers:@[identifier]];
        __block BOOL ok = NO;
        sem = dispatch_semaphore_create(0);
        [center addNotificationRequest:[UNNotificationRequest requestWithIdentifier:identifier content:content trigger:nil] withCompletionHandler:^(NSError *error) { ok = error == nil; dispatch_semaphore_signal(sem); }];
        return spin(sem) && ok;
    }
}
char *pf_notify_language(void) { @autoreleasepool { return strdup((NSLocale.preferredLanguages.firstObject ?: @"en").UTF8String); } }
