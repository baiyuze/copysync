package rpc

import (
	"context"
	"testing"

	"github.com/baiyuze/copysync/client-core/internal/store"
	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
)

func TestListHistorySearch(t *testing.T) {
	var gotQuery string
	s := NewServer(Deps{SearchHistory: func(query string, limit int) ([]store.Clip, error) {
		gotQuery = query
		return []store.Clip{
			{ID: "t", Kind: store.KindText, TextPreview: "季度汇报"},
			{ID: "f", Kind: store.KindFile, Items: []store.Item{{Name: "季度汇报.key"}}},
		}, nil
	}}, NewHub(), "test-token")

	resp, err := s.ListHistory(context.Background(), &pb.ListHistoryRequest{Query: "  季度  ", Limit: 200})
	if err != nil {
		t.Fatal(err)
	}
	// 界面靠 searched 区分「搜过了」与旧版后台服务原样返回的全部记录
	if !resp.GetSearched() || len(resp.GetRecords()) != 2 || gotQuery != "季度" {
		t.Fatalf("searched=%v records=%d query=%q", resp.GetSearched(), len(resp.GetRecords()), gotQuery)
	}

	resp, err = s.ListHistory(context.Background(), &pb.ListHistoryRequest{
		Query: "季度", KindFilter: pb.ClipKind_CLIP_KIND_FILE})
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.GetRecords()) != 1 || resp.GetRecords()[0].GetId() != "f" {
		t.Errorf("搜索结果没有按类型过滤：%v", resp.GetRecords())
	}
}
