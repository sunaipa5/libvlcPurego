//go:build windows

package libvlcPurego

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	procEnumWindows          = user32.NewProc("EnumWindows")
	procGetWindowTextW       = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW = user32.NewProc("GetWindowTextLengthW")
	procSetWinEventHook      = user32.NewProc("SetWinEventHook")
	procGetMessage           = user32.NewProc("GetMessageW")
)

const (
	EVENT_OBJECT_DESTROY  = 0x8001
	WINEVENT_OUTOFCONTEXT = 0x0000
)

type HWND uintptr
type HOOK uintptr
type MSG struct {
	HWnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

func (p *Player) WindowCloseEvent(windowTitle string) <-chan struct{} {
	disptr, winptr, err := findWindowByTitle(windowTitle)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return windowCloseEvent(disptr, winptr)
}

func findWindowByTitle(target string) (HWND, uintptr, error) {
	var foundHWND HWND = 0

	cb := syscall.NewCallback(func(hwnd HWND, lparam uintptr) uintptr {
		length, _, _ := procGetWindowTextLengthW.Call(uintptr(hwnd))
		if length == 0 {
			return 1
		}

		buf := make([]uint16, length+1)
		procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))

		windowTitle := syscall.UTF16ToString(buf)
		if strings.Contains(windowTitle, target) {
			foundHWND = hwnd
			return 0
		}
		return 1
	})

	procEnumWindows.Call(cb, 0)

	if foundHWND == 0 {
		return 0, 0, fmt.Errorf("no window found with title containing: %s", target)
	}

	return foundHWND, 0, nil
}

// in windows not required windowptr, type 0 ptr area
func windowCloseEvent(hwnd HWND, ptr uintptr) <-chan struct{} {
	done := make(chan struct{})

	callback := syscall.NewCallback(func(
		hWinEventHook HOOK,
		event uint32,
		hwnd2 HWND,
		idObject, idChild uint32,
		dwEventThread, dwmsEventTime uint32,
	) uintptr {
		if hwnd2 == hwnd && event == EVENT_OBJECT_DESTROY {
			fmt.Printf("🛑 Window 0x%x was closed.\n", hwnd)
			close(done)
		}
		return 0
	})

	hHook, _, _ := procSetWinEventHook.Call(
		EVENT_OBJECT_DESTROY,
		EVENT_OBJECT_DESTROY,
		0,
		callback,
		0,
		0,
		WINEVENT_OUTOFCONTEXT,
	)

	_ = hHook

	go func() {
		var msg MSG
		for {
			procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		}
	}()

	return done
}
