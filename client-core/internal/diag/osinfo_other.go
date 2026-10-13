//go:build !darwin && !windows && !linux

package diag

import "runtime"

func osVersion() string { return runtime.GOOS }
