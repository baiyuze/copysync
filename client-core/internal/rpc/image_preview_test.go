package rpc

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	pb "github.com/baiyuze/copysync/proto/gen/copysync/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestImagePreviewRPC(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	called := 0
	s := NewServer(Deps{ImagePreviewPath: func(_ context.Context, id string) (string, error) {
		called++
		if id == "missing" {
			return "", errors.New("图片已被清理")
		}
		return "/test/中文 图片.png", nil
	}}, NewHub(), "test-token")
	port, server, err := s.Serve(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer server.Stop()
	conn, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewDaemonServiceClient(conn)
	requestCtx, requestCancel := context.WithTimeout(ctx, 5*time.Second)
	defer requestCancel()
	if _, err := client.GetImagePreview(requestCtx, &pb.GetImagePreviewRequest{ClipId: "image"}); status.Code(err) != codes.Unauthenticated {
		t.Fatalf("unauthenticated: %v", err)
	}
	if called != 0 {
		t.Fatal("unauthenticated call reached image resolver")
	}
	auth := metadata.AppendToOutgoingContext(requestCtx, tokenHeader, "test-token")
	if _, err := client.GetImagePreview(auth, &pb.GetImagePreviewRequest{}); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("empty ID: %v", err)
	}
	if _, err := client.GetImagePreview(auth, &pb.GetImagePreviewRequest{ClipId: "missing"}); status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("missing image: %v", err)
	}
	got, err := client.GetImagePreview(auth, &pb.GetImagePreviewRequest{ClipId: "image"})
	if err != nil || got.GetPath() != "/test/中文 图片.png" {
		t.Fatalf("preview response: %v, %v", got, err)
	}
}
