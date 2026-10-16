package rpc

import (
	"time"

	"github.com/baiyuze/pasteferry/client-core/internal/clipboard"
	"github.com/baiyuze/pasteferry/client-core/internal/config"
	"github.com/baiyuze/pasteferry/client-core/internal/imagepreview"
	"github.com/baiyuze/pasteferry/client-core/internal/store"
	pb "github.com/baiyuze/pasteferry/proto/gen/pasteferry/v1"
)

// 内部类型与 proto 类型的互转。store 层刻意不直接用 proto 类型，
// 以免持久化格式被 wire 格式绑死。

func clipToProto(c store.Clip) *pb.ClipRecord {
	items := make([]*pb.ClipItem, 0, len(c.Items))
	for _, it := range c.Items {
		items = append(items, &pb.ClipItem{
			Name:        it.Name,
			Size:        it.Size,
			IsDir:       it.IsDir,
			ContentType: it.ContentType,
		})
	}
	return &pb.ClipRecord{
		Id:               c.ID,
		Kind:             pb.ClipKind(c.Kind),
		Status:           pb.ClipStatus(c.Status),
		OriginDeviceId:   c.OriginDeviceID,
		OriginDeviceName: c.OriginDeviceName,
		Outgoing:         c.Outgoing,
		Items:            items,
		TotalSize:        c.TotalSize,
		TextPreview:      c.TextPreview,
		CreatedAtUnix:    c.CreatedAt.Unix(),
		ExpiresAtUnix:    c.ExpiresAt.Unix(),
		Error:            c.Error,
		ImageCount:       imagepreview.ImageCount(c),
	}
}

func configToProto(c config.Config) *pb.Config {
	return &pb.Config{
		AutoSyncThresholdBytes: c.AutoSyncThresholdBytes,
		HistoryTtlSeconds:      int64(c.HistoryTTL.Std() / time.Second),
		CacheTtlSeconds:        int64(c.CacheTTL.Std() / time.Second),
		SyncText:               c.SyncText,
		SyncHtml:               c.SyncHTML,
		SyncImage:              c.SyncImage,
		SyncFile:               c.SyncFile,
		AutoApplyToClipboard:   c.AutoApplyToClipboard,
		LaunchAtLogin:          c.LaunchAtLogin,
		DeviceName:             c.DeviceName,
		SignalingUrl:           c.SignalingURL,
		OnlyOwnStun:            c.OnlyOwnSTUN,
		Language:               c.Language,
		LanDisabled:            c.LANDisabled,
		ServiceMode:            serviceModeToProto(c), CustomSignalingUrl: c.CustomSignalingURL, OfficialNoticeAck: c.OfficialNoticeAck,
	}
}

// configFromProto 以 base 打底再覆盖，这样 UI 只改动部分字段时
// 不会把没填的字段清零（proto3 无法区分"未设置"与"零值"）。
func configFromProto(in *pb.Config, base config.Config) config.Config {
	c := base
	if in == nil {
		return c
	}
	if in.GetAutoSyncThresholdBytes() > 0 {
		c.AutoSyncThresholdBytes = in.GetAutoSyncThresholdBytes()
	}
	if in.GetHistoryTtlSeconds() > 0 {
		c.HistoryTTL = config.Duration(time.Duration(in.GetHistoryTtlSeconds()) * time.Second)
	}
	if in.GetCacheTtlSeconds() > 0 {
		c.CacheTTL = config.Duration(time.Duration(in.GetCacheTtlSeconds()) * time.Second)
	}
	if in.GetDeviceName() != "" {
		c.DeviceName = in.GetDeviceName()
	}
	// 布尔项没有"未设置"语义，UI 每次都会带全量值，直接覆盖
	c.SyncText = in.GetSyncText()
	c.SyncHTML = in.GetSyncHtml()
	c.SyncImage = in.GetSyncImage()
	c.SyncFile = in.GetSyncFile()
	c.AutoApplyToClipboard = in.GetAutoApplyToClipboard()
	c.LaunchAtLogin = in.GetLaunchAtLogin()
	c.OnlyOwnSTUN = in.GetOnlyOwnStun()
	// 空字符串就是「跟随系统」，与「没设置」同义，直接覆盖
	c.Language = in.GetLanguage()
	// 服务器地址同理：空字符串就是不用服务器
	c.SignalingURL = in.GetSignalingUrl()
	if in.GetServiceMode() != pb.ServiceMode_SERVICE_MODE_UNSPECIFIED {
		c.CustomSignalingURL = in.GetCustomSignalingUrl()
	}
	switch in.GetServiceMode() {
	case pb.ServiceMode_SERVICE_MODE_OFFICIAL:
		c.ServiceMode = "official"
		c.SignalingURL = ""
	case pb.ServiceMode_SERVICE_MODE_CUSTOM:
		c.ServiceMode = "custom"
		c.CustomSignalingURL = in.GetCustomSignalingUrl()
		if c.CustomSignalingURL == "" {
			c.CustomSignalingURL = in.GetSignalingUrl()
		}
		c.SignalingURL = c.CustomSignalingURL
	case pb.ServiceMode_SERVICE_MODE_LAN:
		c.ServiceMode = "lan"
		c.SignalingURL = ""
	default:
		if c.SignalingURL != "" {
			c.ServiceMode = "custom"
			c.CustomSignalingURL = c.SignalingURL
		} else {
			c.ServiceMode = "lan"
		}
	}
	if in.GetCustomSignalingUrl() != "" {
		c.CustomSignalingURL = in.GetCustomSignalingUrl()
	}
	c.OfficialNoticeAck = in.GetOfficialNoticeAck()
	c.LANDisabled = in.GetLanDisabled()
	return c
}

func serviceModeToProto(c config.Config) pb.ServiceMode {
	switch c.ServiceMode {
	case "official":
		return pb.ServiceMode_SERVICE_MODE_OFFICIAL
	case "custom":
		return pb.ServiceMode_SERVICE_MODE_CUSTOM
	case "lan":
		return pb.ServiceMode_SERVICE_MODE_LAN
	default:
		if c.SignalingURL != "" {
			return pb.ServiceMode_SERVICE_MODE_CUSTOM
		}
		return pb.ServiceMode_SERVICE_MODE_LAN
	}
}

func permissionToProto(p clipboard.Permission) pb.ClipboardPermission {
	switch p {
	case clipboard.PermissionNotApplicable:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_NOT_APPLICABLE
	case clipboard.PermissionDefault:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_DEFAULT
	case clipboard.PermissionAsk:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_ASK
	case clipboard.PermissionAlwaysAllow:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_ALWAYS_ALLOW
	case clipboard.PermissionAlwaysDeny:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_ALWAYS_DENY
	default:
		return pb.ClipboardPermission_CLIPBOARD_PERMISSION_UNSPECIFIED
	}
}

// ClipToProto exposes the same record projection to desktop notifications.
func ClipToProto(c store.Clip) *pb.ClipRecord { return clipToProto(c) }
