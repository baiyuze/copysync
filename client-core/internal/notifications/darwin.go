//go:build darwin

package notifications

/*
#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Cocoa -framework UserNotifications
#include <stdlib.h>
int pf_notify_init(void);
int pf_notify_permission(int request);
int pf_notify_show(const char *id, const char *title, const char *body, const char *action, const char *label);
char *pf_notify_language(void);
*/
import "C"

import (
	"context"
	"errors"
	"github.com/baiyuze/pasteferry/client-core/internal/clipboard"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
	"sync"
	"unsafe"
)

var macAction struct {
	sync.RWMutex
	callback func(string, string)
}

//export goNotificationAction
func goNotificationAction(id, action *C.char) {
	clipID, verb := C.GoString(id), C.GoString(action)
	macAction.RLock()
	callback := macAction.callback
	macAction.RUnlock()
	if callback != nil {
		callback(clipID, verb)
	}
}

type macBackend struct{ loop *clipboard.MainLoop }

func newBackend(loop *clipboard.MainLoop, action func(string, string)) (Backend, error) {
	// New is called on the daemon's locked main thread after InitPlatform.
	if C.pf_notify_init() == 0 {
		return nil, errors.New("notifications require the packaged daemon app")
	}
	macAction.Lock()
	macAction.callback = action
	macAction.Unlock()
	return &macBackend{loop}, nil
}
func (b *macBackend) permission(ctx context.Context, request bool) (*pb.NotificationStatus, error) {
	value := C.int(0)
	if request {
		value = 1
	}
	state, err := clipboard.CallResult(ctx, b.loop, func() (int, error) { return int(C.pf_notify_permission(value)), nil })
	permission := "unavailable"
	switch state {
	case 0:
		permission = "not_requested"
	case 1:
		permission = "denied"
	case 2, 3, 4:
		permission = "allowed"
	}
	return &pb.NotificationStatus{Supported: true, ActionsSupported: true, Permission: permission}, err
}
func (b *macBackend) Status(ctx context.Context) (*pb.NotificationStatus, error) {
	return b.permission(ctx, false)
}
func (b *macBackend) Request(ctx context.Context) (*pb.NotificationStatus, error) {
	return b.permission(ctx, true)
}
func (b *macBackend) Show(ctx context.Context, n Message) error {
	status, err := b.Status(ctx)
	if err != nil || status.Permission != "allowed" {
		return err
	}
	_, err = clipboard.CallResult(ctx, b.loop, func() (bool, error) {
		values := []string{n.ID, n.Title, n.Body, n.Action, n.ActionLabel}
		cs := make([]*C.char, len(values))
		for i, value := range values {
			cs[i] = C.CString(value)
			defer C.free(unsafe.Pointer(cs[i]))
		}
		if C.pf_notify_show(cs[0], cs[1], cs[2], cs[3], cs[4]) == 0 {
			return false, errors.New("notification delivery failed")
		}
		return true, nil
	})
	return err
}
func (*macBackend) Close() { macAction.Lock(); macAction.callback = nil; macAction.Unlock() }
func systemLanguage() string {
	value := C.pf_notify_language()
	if value == nil {
		return "en"
	}
	defer C.free(unsafe.Pointer(value))
	return C.GoString(value)
}
