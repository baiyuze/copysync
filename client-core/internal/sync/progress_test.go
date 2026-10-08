package sync

import (
	"testing"
	"time"

	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

// 进度上报要节流：按数据块上报会把事件队列挤满，挤掉「传输完成」的记录更新。
func TestEmitProgressThrottles(t *testing.T) {
	var got []*pb.TransferProgress
	e := &Engine{
		onProgress:   func(p *pb.TransferProgress) { got = append(got, p) },
		lastProgress: make(map[string]time.Time),
	}
	start := time.Now()
	const total = 10_000
	for sent := int64(100); sent <= total; sent += 100 {
		e.emitProgress("clip", sent, total, start)
	}

	if len(got) > 5 {
		t.Errorf("连续 100 次上报产生了 %d 个事件，应被节流到个位数", len(got))
	}
	// 最后一条必须是 100%，界面靠它知道传完了
	if last := got[len(got)-1]; last.GetTransferred() != total {
		t.Errorf("最后一条进度 = %d，应为 %d", last.GetTransferred(), total)
	}
}
