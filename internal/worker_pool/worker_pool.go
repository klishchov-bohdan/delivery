package worker_pool

import (
	"sync"
)

type Worker interface {
	Stop()
	Do(data interface{}, handlerIndex int)
}

type Constructor func() Worker

type WorkerPool struct {
	WorkerCount int
	DataSource  chan interface{}
	stop        chan struct{}
	New         Constructor
}

func NewPool(count int, new Constructor) *WorkerPool {
	return &WorkerPool{
		WorkerCount: count,
		DataSource:  make(chan interface{}),
		New:         new,
		stop:        make(chan struct{}),
	}
}

func (w *WorkerPool) Stop() {
	for i := 0; i < w.WorkerCount; i++ {
		w.stop <- struct{}{}
	}
}

func (w *WorkerPool) Run() {
	var wg sync.WaitGroup

	for i := 0; i < w.WorkerCount; i++ {
		worker := w.New()
		go func(index int) {
			for {
				select {
				case data, ok := <-w.DataSource:
					if !ok {
						worker.Stop()
						wg.Done()
						return
					}
					if data == nil {
						continue
					}

					worker.Do(data, index)
				case <-w.stop:
					worker.Stop()
					wg.Done()
					return
				}
			}
		}(i)
		wg.Add(1)
	}
	wg.Wait()
}
