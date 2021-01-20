package pool

import (
	"combotest/internal/app/closessl"
	"sync"

	"github.com/sirupsen/logrus"
)

// WorkersPool - пул воркеров, беруших из пула события и что-то делающие с ними
// jobs - откуда брать  работу
// out - куда класть выполненную работу
// workers - число воркеров
type WorkersPool interface {
	Init(jobs EventPool, out EventPool, workers uint8)
	Stop()
}

type encryptWorkersPool struct {
	jobs    EventPool
	done    EventPool
	size    uint8
	wg      *sync.WaitGroup
	running bool
	stop    chan struct{}
	log     *logrus.Logger
}

// NewEncryptWorkersPool -
func NewEncryptWorkersPool(log *logrus.Logger) WorkersPool {
	return &encryptWorkersPool{
		wg:      &sync.WaitGroup{},
		running: false,
		log:     log,
	}
}

func (wp *encryptWorkersPool) Init(jobs EventPool, out EventPool, workers uint8) {
	if wp.running {
		wp.Stop()
	}

	wp.size = workers
	wp.stop = make(chan struct{}, wp.size)
	wp.jobs = jobs
	wp.done = out

	wp.wg.Add(int(wp.size))
	for i := 0; i < int(wp.size); i++ {
		go wp.worker()
	}

}

func (wp *encryptWorkersPool) Stop() {
	for i := 0; i < int(wp.size); i++ {
		wp.stop <- struct{}{}
	}
	wp.wg.Wait()
	wp.running = false
}

func (wp *encryptWorkersPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.stop:
			wp.log.Debug("worker recive stop")
			return // сигнал остановки работы
		default:
			event, ok := wp.jobs.Pop()
			if !ok {
				return // выполнять больше нечего
			}

			data := []byte(event.Message) // copy of event.Message

			event.Message = string(closessl.Encrypt(data))

			wp.done.Push(event)
		}
	}
}
