//go:build windows

package lan

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Windows 上由系统的 DNS-SD 服务（Dnscache）替我们回答 mDNS 查询。
// 实测：后台服务被防火墙挡住入站时，别的电脑照样能通过它发现本机。
const useSystemAdvertiser = true

var (
	dnsapi                   = windows.NewLazySystemDLL("dnsapi.dll")
	procConstructInstance    = dnsapi.NewProc("DnsServiceConstructInstance")
	procFreeInstance         = dnsapi.NewProc("DnsServiceFreeInstance")
	procDnsServiceRegister   = dnsapi.NewProc("DnsServiceRegister")
	procDnsServiceDeRegister = dnsapi.NewProc("DnsServiceDeRegister")

	// 回调只能建有限个且不会释放，全局建一个；注册与注销由 sysMu 串行化，结果经 sysDone 传回
	sysCallback = windows.NewCallback(func(status uint32, _ uintptr, _ uintptr) uintptr {
		select {
		case sysDone <- status:
		default:
		}
		return 0
	})
	sysDone = make(chan uint32, 1)
	sysMu   sync.Mutex
)

const dnsRequestPending = 9506 // DNS_REQUEST_PENDING

// dnsServiceRegisterRequest 对应 DNS_SERVICE_REGISTER_REQUEST。
type dnsServiceRegisterRequest struct {
	Version         uint32
	InterfaceIndex  uint32
	ServiceInstance uintptr
	Callback        uintptr
	Context         uintptr
	Credentials     uintptr
	UnicastEnabled  int32
}

type systemAdvertiser struct {
	req *dnsServiceRegisterRequest
}

func (a *systemAdvertiser) register(name string, port int, text map[string]string) error {
	a.unregister()

	ip, ok := defaultRouteIP()
	if !ok {
		return errors.New("没有可用的局域网地址")
	}
	host, _ := os.Hostname()
	svc, err := windows.UTF16PtrFromString(name + "." + ServiceType + ".local")
	if err != nil {
		return err
	}
	hostName, err := windows.UTF16PtrFromString(host + ".local")
	if err != nil {
		return err
	}
	keys := make([]*uint16, 0, len(text))
	values := make([]*uint16, 0, len(text))
	for _, k := range sortedKeys(text) {
		kp, err1 := windows.UTF16PtrFromString(k)
		vp, err2 := windows.UTF16PtrFromString(text[k])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("TXT 记录含非法字符: %q", k)
		}
		keys = append(keys, kp)
		values = append(values, vp)
	}
	var keysPtr, valuesPtr uintptr
	if len(keys) > 0 {
		keysPtr = uintptr(unsafe.Pointer(&keys[0]))
		valuesPtr = uintptr(unsafe.Pointer(&values[0]))
	}
	// IP4_ADDRESS 是网络字节序的 DWORD
	b := ip.As4()
	addr := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24

	inst, _, callErr := procConstructInstance.Call(
		uintptr(unsafe.Pointer(svc)), uintptr(unsafe.Pointer(hostName)),
		uintptr(unsafe.Pointer(&addr)), 0,
		uintptr(port), 0, 0, uintptr(len(keys)), keysPtr, valuesPtr)
	if inst == 0 {
		return fmt.Errorf("DnsServiceConstructInstance: %v", callErr)
	}

	req := &dnsServiceRegisterRequest{Version: 1, ServiceInstance: inst, Callback: sysCallback}
	if err := sysCall(procDnsServiceRegister, req); err != nil {
		procFreeInstance.Call(inst)
		return fmt.Errorf("DnsServiceRegister: %w", err)
	}
	a.req = req
	return nil
}

func (a *systemAdvertiser) unregister() {
	if a.req == nil {
		return
	}
	req := a.req
	a.req = nil
	_ = sysCall(procDnsServiceDeRegister, req)
	procFreeInstance.Call(req.ServiceInstance)
}

// sysCall 发起一次异步的注册或注销，等系统回调给出结果。
func sysCall(proc *windows.LazyProc, req *dnsServiceRegisterRequest) error {
	if err := proc.Find(); err != nil {
		return err // 系统太旧，没有这组 API
	}
	sysMu.Lock()
	defer sysMu.Unlock()
	select { // 清掉上一次超时后才到的结果
	case <-sysDone:
	default:
	}
	r, _, _ := proc.Call(uintptr(unsafe.Pointer(req)), 0)
	if r != dnsRequestPending {
		return fmt.Errorf("返回 %d", r)
	}
	select {
	case st := <-sysDone:
		if st != 0 {
			return fmt.Errorf("状态 %d", st)
		}
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("等待系统回调超时")
	}
}
