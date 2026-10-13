package diag

import (
	"bufio"
	"os"
	"strings"
)

// osVersion 写发行版的名字，再加上桌面会话是 X11 还是 Wayland：剪贴板的问题常常和它有关
func osVersion() string {
	name := "Linux"
	if f, err := os.Open("/etc/os-release"); err == nil {
		defer f.Close()
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if v, ok := strings.CutPrefix(sc.Text(), "PRETTY_NAME="); ok {
				name = strings.Trim(v, `"`)
				break
			}
		}
	}
	if session := os.Getenv("XDG_SESSION_TYPE"); session != "" {
		name += ", " + session
	}
	return name
}
