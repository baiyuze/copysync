// Package clipfilter identifies clipboard transport artifacts, not user content.
package clipfilter

import "strings"

// RemotePlaceholder matches UU Remote's observed file-clipboard placeholder:
// .uuremote_aeawv followed by a decimal identifier. Match the exact signature,
// not all hidden/empty files or everything in a remote application's directory.
// Accept both path separators because stored paths can originate on either OS.
func RemotePlaceholder(name string) bool {
	name = strings.ReplaceAll(name, `\`, "/")
	name = name[strings.LastIndex(name, "/")+1:]
	suffix, ok := strings.CutPrefix(name, ".uuremote_aeawv")
	if !ok || suffix == "" {
		return false
	}
	for _, c := range suffix {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// UserFiles keeps the ordinary files in mixed selections without modifying the
// caller's slice. Filter before stat/scan so a disappearing placeholder cannot
// become a zero-byte history entry or trigger another remote clipboard write.
func UserFiles(paths []string) []string {
	var out []string
	for _, p := range paths {
		if !RemotePlaceholder(p) {
			out = append(out, p)
		}
	}
	return out
}
