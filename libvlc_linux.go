//go:build linux

package libvlcPurego

import (
	"path"

	"github.com/ebitengine/purego"
)

func initLib(libPath string) (uintptr, error) {
	lib_location := "libvlc.so"
	if libPath == "" {
		libPath = ""
	} else {
		lib_location = path.Join(libPath, "libvlc.so")
	}

	libvlc, err := purego.Dlopen(lib_location, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return 0, err
	}

	return libvlc, nil
}
