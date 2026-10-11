//go:build !darwin && !windows

package lan

import "errors"

// RequestPermission 在其他平台上没有可做的事。
func RequestPermission() error { return errors.ErrUnsupported }

func deniedByOS(error) bool { return false }
