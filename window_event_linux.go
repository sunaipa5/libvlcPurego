//go:build linux

package libvlcPurego

import (
	"fmt"
	"runtime"
	"strings"
	"unsafe"

	"github.com/ebitengine/purego"
)

type XEvent struct {
	xtype int32
	Pad   [23]int64 // Padding to match X11 event size
}

var (
	xlibLib uintptr

	xOpenDisplay       func(displayName *byte) uintptr
	xDefaultRootWindow func(display uintptr) uintptr
	xQueryTree         func(display uintptr, window uintptr, root_return, parent_return *uintptr, children_return *uintptr, nchildren_return *uint32) int
	xFetchName         func(display uintptr, w uintptr, windowName **byte) int
	xFree              func(data uintptr)
	xCloseDisplay      func(display uintptr) int

	xSelectInput        func(display uintptr, window uintptr, mask int64) int
	xNextEvent          func(display uintptr, event *XEvent) int
	destroyNotify             = int32(17)
	structureNotifyMask int64 = 1 << 17
)

func init() {
	var err error
	xlibLib, err = purego.Dlopen("libX11.so", purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic(err)
	}

	purego.RegisterLibFunc(&xOpenDisplay, xlibLib, "XOpenDisplay")
	purego.RegisterLibFunc(&xDefaultRootWindow, xlibLib, "XDefaultRootWindow")
	purego.RegisterLibFunc(&xQueryTree, xlibLib, "XQueryTree")
	purego.RegisterLibFunc(&xFetchName, xlibLib, "XFetchName")
	purego.RegisterLibFunc(&xFree, xlibLib, "XFree")
	purego.RegisterLibFunc(&xCloseDisplay, xlibLib, "XCloseDisplay")
	purego.RegisterLibFunc(&xSelectInput, xlibLib, "XSelectInput")
	purego.RegisterLibFunc(&xNextEvent, xlibLib, "XNextEvent")
}

func (p *Player) WindowCloseEvent(windowTitle string) <-chan struct{} {
	disptr, winptr, err := findWindowByTitle(windowTitle)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return windowCloseEvent(disptr, winptr)
}

func findWindowByTitle(target string) (display uintptr, window uintptr, err error) {
	display = xOpenDisplay(nil)
	if display == 0 {
		return 0, 0, fmt.Errorf("cannot open X display")
	}

	root := xDefaultRootWindow(display)

	var rootRet, parentRet uintptr
	var children uintptr
	var nchildren uint32

	if xQueryTree(display, root, &rootRet, &parentRet, &children, &nchildren) == 0 {
		xCloseDisplay(display)
		return 0, 0, fmt.Errorf("XQueryTree failed")
	}
	defer xFree(children)

	childSlice := unsafe.Slice((*uintptr)(unsafe.Pointer(children)), nchildren)

	for _, win := range childSlice {
		var windowName *byte
		if xFetchName(display, win, &windowName) != 0 && windowName != nil {
			goStr := cStringToGoString(windowName)
			xFree(uintptr(unsafe.Pointer(windowName)))

			if strings.Contains(goStr, target) {
				return display, win, nil
			}
		}
	}

	xCloseDisplay(display)
	return 0, 0, fmt.Errorf("no window found with title: %s", target)
}

func windowCloseEvent(display uintptr, window uintptr) <-chan struct{} {
	if display == 0 || window == 0 {
		return nil
	}

	ch := make(chan struct{})

	xSelectInput(display, window, structureNotifyMask)

	go func() {
		runtime.LockOSThread()
		var event XEvent

		for {
			xNextEvent(display, &event)
			if event.xtype == destroyNotify {
				fmt.Printf("🛑 Window 0x%x was closed.\n", window)
				close(ch)
				return
			}
		}
	}()

	return ch
}

func cStringToGoString(ptr *byte) string {
	if ptr == nil {
		return ""
	}
	var bytes []byte
	base := uintptr(unsafe.Pointer(ptr))
	for i := uintptr(0); ; i++ {
		b := *(*byte)(unsafe.Pointer(base + i))
		if b == 0 {
			break
		}
		bytes = append(bytes, b)
	}
	return string(bytes)
}
