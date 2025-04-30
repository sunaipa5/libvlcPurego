//go:build !linux && !windows

package libvlcPurego

import "fmt"

func (p *Player) WindowCloseEvent(windowTitle string) <-chan struct{} {
	fmt.Println("Window close event not supported on this platform.")
	return nil
}
