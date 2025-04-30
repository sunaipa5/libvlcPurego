package libvlcPurego

import (
	"github.com/ebitengine/purego"
)

var libvlc uintptr

var (
	newVLC                  func(argc int, argv **byte) uintptr
	release                 func(instance uintptr)
	mediaNewPath            func(instance uintptr, path *byte) uintptr
	mediaRelease            func(media uintptr)
	mediaPlayerNewFromMedia func(media uintptr) uintptr
	mediaPlayerPlay         func(player uintptr) int
	mediaPlayerRelease      func(player uintptr)
	mediaNewLocation        func(instance uintptr, location *byte) uintptr
	audioSetVolume          func(player uintptr, volume int32) int
	mediaPlayerEventManager func(player uintptr) uintptr
)

// subtitle
var (
	getSubtitleTrack            func(mp uintptr) int32
	getSubtitleTrackCount       func(mp uintptr) int32
	getSubtitleTrackDescription func(mp uintptr) uintptr
	releaseTrackDescriptionList func(mp uintptr)
	addMediaSlave               func(mp uintptr, slaveType uint32, uri uintptr, selectIt bool) int32
	setSubtitleTrack            func(mp uintptr, trackID int32) int32
)

// event manager
var (
	eventAttach                    func(eventManager uintptr, eventType int, callback uintptr, userData uintptr) int
	eventDetach                    func(eventManager uintptr, eventID int)
	mediaEventManager              func(media uintptr) uintptr
	logEventManager                func() uintptr
	mediaListEventManager          func(mediaList uintptr) uintptr
	mediaListPlayerEventManager    func(mediaListPlayer uintptr) uintptr
	mediaDiscovererEventManager    func(mediaDiscoverer uintptr) uintptr
	rendererDiscovererEventManager func(rendererDiscoverer uintptr) uintptr
	retainEvent                    func(event uintptr)
	releaseEvent                   func(event uintptr)
	eventTypeName                  func(eventType int) *byte
)

var (
	getNSObject func(player uintptr) uintptr
	getHWND     func(player uintptr) uintptr
	getXWindow  func(player uintptr) uint32
)

func register_libvlc_funcs() {
	purego.RegisterLibFunc(&mediaNewLocation, libvlc, "libvlc_media_new_location")
	purego.RegisterLibFunc(&newVLC, libvlc, "libvlc_new")
	purego.RegisterLibFunc(&release, libvlc, "libvlc_release")

	purego.RegisterLibFunc(&mediaNewPath, libvlc, "libvlc_media_new_path")
	purego.RegisterLibFunc(&mediaRelease, libvlc, "libvlc_media_release")
	purego.RegisterLibFunc(&mediaPlayerNewFromMedia, libvlc, "libvlc_media_player_new_from_media")
	purego.RegisterLibFunc(&mediaPlayerPlay, libvlc, "libvlc_media_player_play")
	purego.RegisterLibFunc(&mediaPlayerRelease, libvlc, "libvlc_media_player_release")
	purego.RegisterLibFunc(&audioSetVolume, libvlc, "libvlc_audio_set_volume")

	//Subtitle
	purego.RegisterLibFunc(&getSubtitleTrack, libvlc, "libvlc_video_get_spu")
	purego.RegisterLibFunc(&getSubtitleTrackCount, libvlc, "libvlc_video_get_spu_count")
	purego.RegisterLibFunc(&getSubtitleTrackDescription, libvlc, "libvlc_video_get_spu_description")
	purego.RegisterLibFunc(&releaseTrackDescriptionList, libvlc, "libvlc_track_description_list_release")
	purego.RegisterLibFunc(&addMediaSlave, libvlc, "libvlc_media_player_add_slave")
	purego.RegisterLibFunc(&setSubtitleTrack, libvlc, "libvlc_video_set_spu")

	// Event Manager
	purego.RegisterLibFunc(&mediaPlayerEventManager, libvlc, "libvlc_media_player_event_manager")
	purego.RegisterLibFunc(&eventAttach, libvlc, "libvlc_event_attach")
	purego.RegisterLibFunc(&eventDetach, libvlc, "libvlc_event_detach")
	purego.RegisterLibFunc(&mediaEventManager, libvlc, "libvlc_media_event_manager")
	purego.RegisterLibFunc(&logEventManager, libvlc, "libvlc_log_get_context")
	purego.RegisterLibFunc(&mediaListEventManager, libvlc, "libvlc_media_list_event_manager")
	purego.RegisterLibFunc(&mediaListPlayerEventManager, libvlc, "libvlc_media_list_player_event_manager")
	purego.RegisterLibFunc(&mediaDiscovererEventManager, libvlc, "libvlc_media_discoverer_event_manager")
	purego.RegisterLibFunc(&rendererDiscovererEventManager, libvlc, "libvlc_renderer_discoverer_event_manager")
	purego.RegisterLibFunc(&retainEvent, libvlc, "libvlc_event_attach")
	purego.RegisterLibFunc(&releaseEvent, libvlc, "libvlc_event_detach")
	purego.RegisterLibFunc(&eventTypeName, libvlc, "libvlc_event_type_name")

	purego.RegisterLibFunc(&getXWindow, libvlc, "libvlc_media_player_get_xwindow")
	purego.RegisterLibFunc(&getHWND, libvlc, "libvlc_media_player_get_hwnd")
	purego.RegisterLibFunc(&getNSObject, libvlc, "libvlc_media_player_get_nsobject")
}
