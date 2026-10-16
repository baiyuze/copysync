//go:build !darwin && !windows && !linux

package notifications

import "github.com/baiyuze/pasteferry/client-core/internal/clipboard"

func newBackend(*clipboard.MainLoop, func(string, string)) (Backend, error) {
	return unavailable{}, nil
}
func systemLanguage() string { return "en" }
