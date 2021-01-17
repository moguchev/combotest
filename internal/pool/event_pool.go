package pool

import (
	"combotest/internal/models"
)

// EventPool - пул принятых событий от агентов
type EventPool interface {
	// Push - положить событие в пул
	Push(e models.Event)
	// Pop - достать событие из пула
	Pop() (models.Event, bool)

	Close()
}

type chanEventPool struct {
	queue chan models.Event
}

// NewChanEventPool return implementation of EventPool interface
func NewChanEventPool(size uint32) EventPool {
	return &chanEventPool{queue: make(chan models.Event, size)}
}

func (ep *chanEventPool) Push(e models.Event) {
	ep.queue <- e
}

func (ep *chanEventPool) Pop() (models.Event, bool) {
	e, ok := <-ep.queue
	return e, ok
}

func (ep *chanEventPool) Close() {
	close(ep.queue)
}
