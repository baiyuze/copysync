package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/config"
	"github.com/baiyuze/pasteferry/client-core/internal/notifications"
	"github.com/baiyuze/pasteferry/client-core/internal/rpc"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// A Windows protocol activation is a client of the existing daemon, never a
// second clipboard owner or supervisor. The endpoint token stays off the URI.
func activateNotification(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "pasteferry-notify" || u.Host != "action" {
		return fmt.Errorf("invalid notification URI")
	}
	verb, id := u.Query().Get("action"), u.Query().Get("clip")
	if verb == "open" {
		return notifications.OpenUI()
	}
	if (verb != "receive" && verb != "apply") || id == "" {
		return fmt.Errorf("invalid notification action")
	}
	paths, err := config.DefaultPaths()
	if err != nil {
		return err
	}
	ep, err := rpc.ReadEndpoint(paths.Endpoint)
	if err != nil {
		return err
	}
	conn, err := grpc.NewClient(net.JoinHostPort("127.0.0.1", fmt.Sprint(ep.Port)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-pasteferry-token", ep.Token)
	_, err = pb.NewDaemonServiceClient(conn).ActivateNotification(ctx, &pb.NotificationAction{ClipId: id, Action: verb})
	return err
}
