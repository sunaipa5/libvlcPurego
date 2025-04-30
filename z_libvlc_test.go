package libvlcPurego

import (
	"fmt"

	"testing"
)

func Test_player(t *testing.T) {
	if err := Init(); err != nil {
		panic(err)
	}

	player, err := NewPlayer(nil)
	if err != nil {
		t.Errorf("Failed to create player: %v", err)
	}
	defer player.Release()

	if err := player.NewSource("https://sample-videos.com/video321/mp4/720/big_buck_bunny_720p_2mb.mp4"); err != nil {
		t.Errorf("Failed to set source: %v", err)
	}

	player.Play()

	eventManager, err := NewEventManager(player.Player)
	if err != nil {
		panic(err)
	}
	fmt.Println("Event Manager Created")

	vout := make(chan struct{})
	eventid, err := eventManager.EventListenerOld(MediaPlayerVout, vout)
	if err != nil {
		panic(err)
	}
	fmt.Println(eventid)

	<-vout
	fmt.Println("Vout Event Recivet")

	//windows: VLC (Direct3D11 output)
	//Linux:  VLC Media Player
	closeChan := player.WindowCloseEvent("VLC Media Player")
	<-closeChan
	fmt.Println("Player window closed")
	player.Release()
	fmt.Println("Player released")

}
