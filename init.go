package libvlcPurego

import (
	"fmt"
	"runtime"
)

type CustomLibPaths struct {
	Windows string
	Linux   string
	Darwin  string
}

// Init with default locations
func Init() error {
	var err error
	libvlc, err = initLib("")
	if err != nil {
		return err
	}

	register_libvlc_funcs()

	return nil
}

/*
Init with custom library paths.
Default paths are used for empty fields.

# Example

	customPaths := libvlcPurego.CustomLibPaths{
		Windows: "C:/Program Files/VideoLAN/VLC",
		Linux:   "/usr/lib",
		Darwin:  "/usr/lib",
	}

	if err := customPaths.InitCustom(); err != nil {
		panic(err)
	}
*/
func (p CustomLibPaths) InitCustom() error {

	var err error
	switch runtime.GOOS {
	case "linux":
		libvlc, err = initLib(p.Linux)
		if err != nil {
			return err
		}
		register_libvlc_funcs()
	case "windows":
		libvlc, err = initLib(p.Linux)
		if err != nil {
			return err
		}
		register_libvlc_funcs()
	case "darwin":
		libvlc, err = initLib(p.Linux)
		if err != nil {
			return err
		}
		register_libvlc_funcs()
	default:
		return fmt.Errorf("unsupported platform")
	}

	return nil
}
