package sync

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
)

// chunkWriter 把连续字节流切成 DataChannel 能承受的小块并施加背压。
//
// 两件事缺一不可：
//  1. 分片——SCTP 对超大消息的支持在各实现间差异很大，60 KiB 是稳妥值；
//  2. 背压——不管发送缓冲直接灌数据，传大文件会把内存吃光。
//     缓冲超过高水位时暂停读取源文件，等它回落再继续。
type chunkWriter struct {
	dc         *webrtc.DataChannel
	onProgress func(sent int64)

	buf  []byte
	sent int64

	mu  sync.Mutex
	err error
}

func newChunkWriter(dc *webrtc.DataChannel, onProgress func(int64)) *chunkWriter {
	w := &chunkWriter{
		dc:         dc,
		onProgress: onProgress,
		buf:        make([]byte, 0, chunkSize),
	}
	dc.OnError(func(err error) {
		w.mu.Lock()
		w.err = err
		w.mu.Unlock()
	})
	return w
}

func (w *chunkWriter) Write(p []byte) (int, error) {
	total := len(p)
	for len(p) > 0 {
		space := chunkSize - len(w.buf)
		n := min(space, len(p))
		w.buf = append(w.buf, p[:n]...)
		p = p[n:]

		if len(w.buf) >= chunkSize {
			if err := w.flushChunk(); err != nil {
				return total - len(p), err
			}
		}
	}
	return total, nil
}

func (w *chunkWriter) Flush() error {
	if len(w.buf) == 0 {
		return nil
	}
	return w.flushChunk()
}

func (w *chunkWriter) flushChunk() error {
	if err := w.waitForBuffer(); err != nil {
		return err
	}
	if w.dc.ReadyState() != webrtc.DataChannelStateOpen {
		return errors.New("数据通道已关闭")
	}
	if err := w.dc.Send(w.buf); err != nil {
		return err
	}
	w.sent += int64(len(w.buf))
	w.buf = w.buf[:0]

	if w.onProgress != nil {
		w.onProgress(w.sent)
	}
	return nil
}

// waitForBuffer 在发送缓冲积压过多时等待它回落。
//
// 刻意用轮询而非 OnBufferedAmountLow 回调：那个回调是边沿触发的，
// 若缓冲在等待方进入 Wait 之前就已降到阈值以下，回调早已发生过，
// 等待就再也等不到唤醒——传 20 MB 文件时实测会卡死在 786 KB 处。
// 2 ms 的轮询间隔对文件传输而言开销可忽略，却彻底消除了丢失唤醒。
func (w *chunkWriter) waitForBuffer() error {
	const (
		pollInterval = 2 * time.Millisecond
		maxWait      = 2 * time.Minute
	)
	deadline := time.Now().Add(maxWait)

	for w.dc.BufferedAmount() >= bufferHighWater {
		w.mu.Lock()
		err := w.err
		w.mu.Unlock()
		if err != nil {
			return err
		}
		if w.dc.ReadyState() != webrtc.DataChannelStateOpen {
			return errors.New("数据通道在传输中被关闭")
		}
		if time.Now().After(deadline) {
			return errors.New("等待发送缓冲排空超时")
		}
		time.Sleep(pollInterval)
	}
	return nil
}

// WaitDrained 等待发送缓冲彻底排空。
//
// 数据通道一关，尚未发出的数据就丢了；关闭前必须确认本地缓冲已清空。
func (w *chunkWriter) WaitDrained(ctx context.Context, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if w.dc.BufferedAmount() == 0 {
			return
		}
		if w.dc.ReadyState() != webrtc.DataChannelStateOpen {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Millisecond):
		}
	}
}

func (w *chunkWriter) Sent() int64 { return w.sent }
