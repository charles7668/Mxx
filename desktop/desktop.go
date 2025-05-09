package desktop

import (
	"fmt"
	"github.com/pkg/browser"
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

func isDoubleClickLaunched() bool {
	lp := kernel32.NewProc("GetConsoleProcessList")
	if lp != nil {
		var pids [2]uint32
		var maxCount uint32 = 2
		ret, _, _ := lp.Call(uintptr(unsafe.Pointer(&pids)), uintptr(maxCount))
		if ret > 1 {
			return false
		}
	}
	return true
}

func openURLInBrowser(hostAddress string) {
	err := browser.OpenURL(hostAddress)
	if err != nil {
		_ = fmt.Errorf("Could not open browser: %s\n", err.Error())
	}
}

func Launch(address string) {
	openURLInBrowser(address)
}

func IsDesktop() bool {
	if isDoubleClickLaunched() {
		return true
	}

	return false
}
