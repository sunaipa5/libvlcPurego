//go:build windows

package libvlcPurego

import (
	"os"

	"syscall"
)

func initLib(libPath string) (uintptr, error) {
	if libPath == "" {
		libPath = "C:/Program Files/VideoLAN/VLC"
	}

	os.Setenv("PATH", os.Getenv("PATH")+";"+libPath)

	libvlc, err := syscall.LoadLibrary("libvlc.dll")
	if err != nil {
		return 0, err
	}

	return uintptr(libvlc), nil
}
