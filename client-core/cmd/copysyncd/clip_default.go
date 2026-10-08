//go:build !fakeclip

package main

import "github.com/baiyuze/copysync/client-core/internal/clipboard"

func newClipboard() clipboard.Clipboard { return clipboard.New() }
