package libvlcPurego

import (
	"fmt"
	"sync"

	"github.com/ebitengine/purego"
)

type EventManager struct {
	manager      uintptr
	eventStore   map[uintptr]chan struct{}
	eventCounter uintptr
	eventLock    sync.Mutex
}

func NewEventManager(player uintptr) (EventManager, error) {
	eventManager := mediaPlayerEventManager(player)
	if eventManager == 0 {
		return EventManager{}, fmt.Errorf("failed to get event manager")
	}

	return EventManager{
		manager:      eventManager,
		eventStore:   make(map[uintptr]chan struct{}),
		eventCounter: 1,
	}, nil
}

func (em *EventManager) EventListenerOld(event Event, ch chan struct{}) (int, error) {
	callback := purego.NewCallback(func(event uintptr, userCh uintptr) uintptr {
		em.eventLock.Lock()
		ch, ok := em.eventStore[userCh]
		if ok && ch != nil {
			close(ch)
			delete(em.eventStore, userCh)
		}
		em.eventLock.Unlock()

		return 0
	})

	em.eventLock.Lock()
	id := em.eventCounter
	em.eventCounter++
	em.eventStore[id] = ch
	em.eventLock.Unlock()

	eventid := eventAttach(em.manager, int(event), callback, id)
	if eventid != 0 {
		return 0, fmt.Errorf("Event attach failed")
	}

	return eventid, nil
}
