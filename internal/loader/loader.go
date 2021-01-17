package loader

import (
	"combotest/internal/events"
	"combotest/internal/models"
	"combotest/internal/pool"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type Loader struct {
	dataPool pool.EventPool
	storage  events.Repository
	wg       *sync.WaitGroup
	stop     chan struct{}
	loaders  int
	log      *logrus.Logger
}

func NewLoader(dp pool.EventPool, st events.Repository, log *logrus.Logger) *Loader {
	return &Loader{
		dataPool: dp,
		storage:  st,
		log:      log,
	}
}

func (l *Loader) Run(numParallelLoaders, chunkSize uint32) {
	l.loaders = int(numParallelLoaders)
	l.stop = make(chan struct{}, l.loaders)
	l.wg = &sync.WaitGroup{}
	l.wg.Add(l.loaders)

	for i := 0; i < l.loaders; i++ {
		go l.loading(chunkSize)
	}
}

func (l *Loader) Stop() {
	for i := 0; i < l.loaders; i++ {
		l.stop <- struct{}{}
	}
	l.wg.Wait()
}

func (l *Loader) loading(chunkSize uint32) {
	log := l.log.WithField("actor", "loader")
	defer l.wg.Done()

	chunk := make([]models.Event, 0, chunkSize)
load_loop:
	for {
		ch := make(chan models.Event, 1)
		go func() {
			e, ok := l.dataPool.Pop()
			if !ok {
				close(ch)
				return
			}
			ch <- e
		}()

		select {
		case <-l.stop:
			break load_loop
		case e, ok := <-ch:
			if !ok {
				break load_loop
			}

			chunk = append(chunk, e)

			if (len(chunk)) == int(chunkSize) { // отправляем чанк
				if err := l.storage.InsertEvents(context.TODO(), chunk); err != nil {
					log.WithError(err).Error("insert")
				}

				chunk = make([]models.Event, 0, chunkSize)
			}
		}
	}
	// отправляем что осталось
	if len(chunk) > 0 {
		if err := l.storage.InsertEvents(context.TODO(), chunk); err != nil {
			log.WithError(err).Error("insert last")
		}
	}
}
