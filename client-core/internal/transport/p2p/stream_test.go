package p2p_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"testing"
	"time"

	"github.com/baiyuze/copysync/client-core/internal/transport/p2p"
)

// TestStreamDeliversPayloadIntact 是文件传输的核心保障。
//
// 它一次性锁住三个曾经真实发生过的故障：
//  1. 时序：OnMessage 挂得比首个分片晚 → 接收方收到 0 字节；
//  2. 顺序：缓冲回放期间新消息插队 → zstd 流损坏，报 "compressed size too big"；
//  3. 背压：依赖边沿触发的 OnBufferedAmountLow → 丢失唤醒，20 MB 传到 786 KB 卡死。
//
// 随机负载 + 校验和比对能同时覆盖"有没有丢"和"顺序对不对"。
func TestStreamDeliversPayloadIntact(t *testing.T) {
	const payloadSize = 4 << 20 // 4 MiB，足以跨越多个分片与缓冲水位

	r := newRelay()
	alice := newHarness(t, "alice", r)
	bob := newHarness(t, "bob", r)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := alice.manager.Connect(ctx, "bob"); err != nil {
		t.Fatalf("建立连接: %v", err)
	}
	waitConnected(t, alice, "bob")
	waitConnected(t, bob, "alice")

	payload := make([]byte, payloadSize)
	if _, err := rand.Read(payload); err != nil {
		t.Fatal(err)
	}
	want := sha256.Sum256(payload)

	dc, err := alice.manager.OpenStream("bob", "payload")
	if err != nil {
		t.Fatalf("开数据流: %v", err)
	}

	opened := make(chan struct{})
	dc.OnOpen(func() { close(opened) })
	select {
	case <-opened:
	case <-time.After(20 * time.Second):
		t.Fatal("数据流未能打开")
	}

	// 接收方：从队列按序取出并拼回完整内容
	var stream *p2p.Stream
	select {
	case stream = <-bob.streams:
	case <-time.After(20 * time.Second):
		t.Fatal("接收方未收到数据流")
	}

	received := make(chan []byte, 1)
	go func() {
		var buf bytes.Buffer
		for chunk := range stream.Chunks() {
			buf.Write(chunk)
		}
		received <- buf.Bytes()
	}()

	// 发送方：按 60 KiB 分片发出，中途不做任何节流之外的处理
	const chunk = 60 * 1024
	for off := 0; off < len(payload); off += chunk {
		end := min(off+chunk, len(payload))
		// 简单节流，避免测试里把发送缓冲撑爆（生产路径由 chunkWriter 负责背压）
		for dc.BufferedAmount() > 1<<20 {
			time.Sleep(2 * time.Millisecond)
		}
		if err := dc.Send(payload[off:end]); err != nil {
			t.Fatalf("发送分片 @%d: %v", off, err)
		}
	}
	for dc.BufferedAmount() > 0 {
		time.Sleep(5 * time.Millisecond)
	}
	dc.Close()

	select {
	case got := <-received:
		if len(got) != payloadSize {
			t.Fatalf("收到 %d 字节，期望 %d（有分片丢失）", len(got), payloadSize)
		}
		if gotSum := sha256.Sum256(got); gotSum != want {
			t.Error("校验和不一致：分片顺序被打乱了")
		}
	case <-time.After(60 * time.Second):
		t.Fatal("接收超时")
	}
}

// TestStreamBuffersEarlyMessages 验证在消费者挂上来之前到达的分片不会丢失。
//
// pion 的 OnMessage 不缓冲：handler 设置前到达的消息会被直接丢弃。
// Stream 必须在 OnDataChannel 回调里当场接住它们。
func TestStreamBuffersEarlyMessages(t *testing.T) {
	r := newRelay()
	alice := newHarness(t, "alice", r)
	bob := newHarness(t, "bob", r)

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	if err := alice.manager.Connect(ctx, "bob"); err != nil {
		t.Fatalf("建立连接: %v", err)
	}
	waitConnected(t, alice, "bob")
	waitConnected(t, bob, "alice")

	dc, err := alice.manager.OpenStream("bob", "early")
	if err != nil {
		t.Fatalf("开数据流: %v", err)
	}
	opened := make(chan struct{})
	dc.OnOpen(func() { close(opened) })
	select {
	case <-opened:
	case <-time.After(20 * time.Second):
		t.Fatal("数据流未能打开")
	}

	// 立刻发送，不给接收侧任何准备时间
	for i := range 5 {
		if err := dc.Send([]byte{byte('A' + i)}); err != nil {
			t.Fatalf("发送: %v", err)
		}
	}

	var stream *p2p.Stream
	select {
	case stream = <-bob.streams:
	case <-time.After(20 * time.Second):
		t.Fatal("接收方未收到数据流")
	}

	// 故意拖一会儿再消费，模拟上层查表建管道的耗时
	time.Sleep(300 * time.Millisecond)
	dc.Close()

	var got []byte
	timeout := time.After(20 * time.Second)
	for {
		select {
		case chunk, ok := <-stream.Chunks():
			if !ok {
				if string(got) != "ABCDE" {
					t.Errorf("收到 %q，期望 ABCDE（早到的分片被丢弃或乱序）", got)
				}
				return
			}
			got = append(got, chunk...)
		case <-timeout:
			t.Fatalf("接收超时，已收到 %q", got)
		}
	}
}
