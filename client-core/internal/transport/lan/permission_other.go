//go:build !darwin && !windows

package lan

import (
	"errors"
	"net/netip"
)

// RequestPermission 在其他平台上没有可做的事。
func RequestPermission() error { return errors.ErrUnsupported }

func deniedByOS(error) bool { return false }

func probeTarget() (netip.Addr, bool) { return netip.Addr{}, false }
