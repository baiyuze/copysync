package gateway

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/host"
)

var machineOnce sync.Once
var machineDigest string

// MachineDigest is called only from relay authorization, never at startup.
// An unavailable identifier leaves the server's previous binding intact.
func MachineDigest() string {
	machineOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		raw, err := readMachineID(ctx, runtime.GOOS, os.ReadFile, host.HostIDWithContext)
		if err != nil {
			return
		}
		normalized, err := normalizeMachineID(raw)
		if err != nil {
			return
		}
		mac := hmac.New(sha256.New, []byte("PasteFerry official relay quota v1"))
		mac.Write([]byte(normalized))
		machineDigest = "v1:" + hex.EncodeToString(mac.Sum(nil))
	})
	return machineDigest
}

func readMachineID(ctx context.Context, platform string, read func(string) ([]byte, error), hostID func(context.Context) (string, error)) (string, error) {
	if platform != "linux" {
		return hostID(ctx)
	}
	for _, path := range []string{"/etc/machine-id", "/var/lib/dbus/machine-id"} {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		if b, err := read(path); err == nil {
			if _, err = normalizeMachineID(string(b)); err == nil {
				return string(b), nil
			}
		}
	}
	return "", errors.New("stable machine ID unavailable")
}

func normalizeMachineID(raw string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(raw))
	s = strings.ReplaceAll(s, "-", "")
	if len(s) != 32 {
		return "", errors.New("invalid machine ID")
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}
	allZero, allFF := true, true
	for _, v := range b {
		allZero = allZero && v == 0
		allFF = allFF && v == 255
	}
	if allZero || allFF {
		return "", errors.New("invalid machine ID")
	}
	return s, nil
}
