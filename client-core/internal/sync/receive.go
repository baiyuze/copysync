package sync

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/store"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
)

type ReceiveStage int

const (
	ReceivePending ReceiveStage = iota
	ReceiveStarted
	ReceiveCompleted
	ReceiveFailed
)

// ReceiveEvent is live intent, unlike a replayable history record. Preview
// downloads and reconnect backlogs must not produce actionable notifications.
type ReceiveEvent struct {
	Stage     ReceiveStage
	Record    *pb.ClipRecord
	Clipboard *pb.ClipboardWriteResult
}

type receiveAttempt struct {
	manual, preview bool
	revision        *int64
	streamID        string
	timer           *time.Timer
}

func (e *Engine) receiveEvent(clip store.Clip, stage ReceiveStage, result *pb.ClipboardWriteResult) {
	if e.onReceive != nil {
		e.onReceive(ReceiveEvent{stage, recordProto(clip), result})
	}
}

// All receiveAttempt fields and timers are protected by e.mu.
func (e *Engine) newReceiveLocked(id string, manual, preview bool) *receiveAttempt {
	if e.receives == nil {
		e.receives = make(map[string]*receiveAttempt)
	}
	a := &receiveAttempt{manual: manual, preview: preview}
	e.receives[id] = a
	d := e.receiveStartTimeout
	if d <= 0 {
		d = 30 * time.Second
	}
	a.timer = time.AfterFunc(d, func() { e.failReceive(id, a, errors.New("接收超时，请确认源设备在线后重试")) })
	return a
}

func (e *Engine) startAutomaticReceive(clip store.Clip, quiet bool) {
	e.mu.Lock()
	if e.receives[clip.ID] == nil {
		e.newReceiveLocked(clip.ID, false, false)
	}
	if quiet {
		if e.quiet == nil {
			e.quiet = make(map[string]bool)
		}
		e.quiet[clip.ID] = true
	}
	_ = e.store.SetClipStatus(clip.ID, store.StatusFetching, "")
	clip.Status = store.StatusFetching
	e.mu.Unlock()
	e.emitRecord(clip)
}

func (e *Engine) fetch(ctx context.Context, clipID string, preview bool) error {
	e.mu.Lock()
	clip, err := e.store.GetClip(clipID)
	if err != nil {
		e.mu.Unlock()
		return err
	}
	if clip.Outgoing {
		e.mu.Unlock()
		return errors.New("这是本机复制的内容，无需接收")
	}
	if clip.Status == store.StatusReady || e.receives[clipID] != nil {
		e.mu.Unlock()
		return nil
	}
	// Also recovers FETCHING records left behind by a previous daemon process.
	a := e.newReceiveLocked(clipID, true, preview)
	if preview {
		if e.quiet == nil {
			e.quiet = make(map[string]bool)
		}
		e.quiet[clipID] = true
	}
	e.mu.Unlock()
	if !preview && e.watcher != nil {
		if revision, readErr := e.watcher.Revision(ctx); readErr == nil {
			e.mu.Lock()
			a.revision = &revision
			e.mu.Unlock()
		}
	}
	e.mu.Lock()
	if e.receives[clipID] != a {
		e.mu.Unlock()
		return errors.New("接收请求已超时，请重试")
	}
	err = e.store.SetClipStatus(clipID, store.StatusFetching, "")
	e.mu.Unlock()
	if err != nil {
		e.failReceive(clipID, a, err)
		return err
	}
	clip.Status, clip.Error = store.StatusFetching, ""
	e.emitRecord(clip)
	if !preview {
		e.receiveEvent(clip, ReceiveStarted, nil)
	}
	if err := e.send(clip.OriginDeviceID, &pb.PeerMessage{
		Payload: &pb.PeerMessage_Fetch{Fetch: &pb.PeerFetch{ClipId: clipID}},
	}); err != nil {
		wrapped := fmt.Errorf("源设备当前不可达: %w", err)
		e.failReceive(clipID, a, wrapped)
		return wrapped
	}
	return nil
}

func (e *Engine) retireReceiveLocked(id string, a *receiveAttempt) {
	a.timer.Stop()
	delete(e.receives, id)
	delete(e.incoming, a.streamID)
	delete(e.quiet, id)
}

func (e *Engine) failReceive(id string, expected *receiveAttempt, failure error) {
	e.mu.Lock()
	a := e.receives[id]
	if a == nil || (expected != nil && expected != a) {
		e.mu.Unlock()
		return
	}
	r := e.incoming[a.streamID]
	_ = e.store.SetClipStatus(id, store.StatusFailed, failure.Error())
	e.retireReceiveLocked(id, a)
	e.mu.Unlock()
	if r != nil {
		r.pipeW.CloseWithError(failure)
		if r.closeStream != nil {
			r.closeStream()
		}
	}
	e.forgetProgress(id)
	if clip, err := e.store.GetClip(id); err == nil {
		e.emitRecord(clip)
		if a.manual && !a.preview {
			e.receiveEvent(clip, ReceiveFailed, nil)
		}
	}
}

func (e *Engine) writeClipboard(ctx context.Context, clip store.Clip, revision *int64) *pb.ClipboardWriteResult {
	result := &pb.ClipboardWriteResult{ClipId: clip.ID, Outcome: pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_FAILED}
	content, err := e.clipboardContent(clip)
	if err == nil {
		if e.watcher == nil {
			err = errors.New("剪贴板服务不可用")
		} else if revision != nil {
			var applied bool
			applied, err = e.watcher.WriteIfUnchanged(ctx, content, *revision)
			if err == nil && !applied {
				result.Outcome = pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_PRESERVED
				return result
			}
		} else {
			err = e.watcher.Write(ctx, content)
		}
	}
	if err != nil {
		result.Error = err.Error()
	} else {
		result.Outcome = pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_APPLIED
	}
	return result
}

func (e *Engine) clipboardResult(result *pb.ClipboardWriteResult) {
	message := ""
	if result.GetOutcome() == pb.ClipboardWriteOutcome_CLIPBOARD_WRITE_OUTCOME_FAILED {
		message = result.GetError()
	}
	if e.store != nil {
		changed, err := e.store.SetClipboardError(result.GetClipId(), message)
		if err != nil {
			e.log.Warn("保存剪贴板写入结果失败", "err", err)
		} else if changed {
			if clip, err := e.store.GetClip(result.GetClipId()); err == nil {
				e.emitRecord(clip)
			}
		}
	}
	if e.onClipboardResult != nil {
		e.onClipboardResult(result)
	}
}

// Close cancels watchdogs before the store is closed on shutdown.
func (e *Engine) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	for id, a := range e.receives {
		if r := e.incoming[a.streamID]; r != nil {
			r.pipeW.CloseWithError(context.Canceled)
		}
		e.retireReceiveLocked(id, a)
	}
}
