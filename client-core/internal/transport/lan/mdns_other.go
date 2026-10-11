//go:build !windows

package lan

// 其他平台由 pion/mdns 自己回答查询。
const useSystemAdvertiser = false

type systemAdvertiser struct{}

func (*systemAdvertiser) register(string, int, map[string]string) error { return nil }
func (*systemAdvertiser) unregister()                                   {}
