package libvlcPurego

import (
	"fmt"
	"net/url"
	"os"
)

type Player struct {
	Instance uintptr
	Player   uintptr
}

func NewPlayer(args ...string) (*Player, error) {
	var argv **byte = nil
	var argc int = 0

	if len(args) != 0 {
		argc = len(args)
		var cleanup func()
		argv, cleanup = stringSliceToPtrPtrByte(args)
		defer cleanup()
	}

	instance := newVLC(argc, argv)
	if instance == 0 {
		return nil, fmt.Errorf("libvlc_new failed")
	}

	return &Player{
		Instance: instance,
	}, nil
}

// File path or source URL
func (p *Player) NewSource(source string) error {
	u, err := url.Parse(source)
	isURL := err == nil && u.Scheme != ""

	if isURL {
		urlBytes := append([]byte(source), 0)
		media := mediaNewLocation(p.Instance, &urlBytes[0])
		if media == 0 {
			return fmt.Errorf("libvlc_media_new_location failed")
		}

		playermedia := mediaPlayerNewFromMedia(media)
		if playermedia == 0 {
			return fmt.Errorf("libvlc_media_player_new_from_media failed")
		}

		p.Player = playermedia

		return nil
	}

	_, statErr := os.Stat(source)
	isFile := statErr == nil

	if isFile {
		pathBytes := append([]byte(source), 0)
		media := mediaNewPath(p.Instance, &pathBytes[0])
		if media == 0 {
			return fmt.Errorf("libvlc_media_new_paths failed")
		}

		playermedia := mediaPlayerNewFromMedia(media)
		if playermedia == 0 {
			return fmt.Errorf("libvlc_media_player_new_from_media failed")
		}

		p.Player = playermedia
		return nil
	}

	return fmt.Errorf("unsupported source type")
}

func (p *Player) Release() {
	if p.Player != 0 {
		mediaRelease(p.Player)
		p.Player = 0
	}
	if p.Instance != 0 {
		release(p.Instance)
		p.Instance = 0
	}
}

func (p *Player) Play() {
	result := mediaPlayerPlay(p.Player)
	if result != 0 {
		panic("libvlc_media_player_play failed")
	}
}
