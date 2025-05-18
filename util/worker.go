package util

import (
	"log/slog"
	"sync"
)

type WorkerPool struct {
	taskChan chan func() error
	wg       sync.WaitGroup
	once     sync.Once
}

var threadPool = getWorkerPool(4)

func getWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		taskChan: make(chan func() error, 100),
	}
	pool.wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go pool.worker()
	}
	return pool
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()
	for task := range p.taskChan {
		if task != nil {
			if err := task(); err != nil {
				slog.Error("Function task failed: " + err.Error())
			}
		}
	}
}

func (p *WorkerPool) AddTask(task func() error) {
	p.taskChan <- task
}

func (p *WorkerPool) Close() {
	p.once.Do(func() {
		close(p.taskChan)
		p.wg.Wait()
	})
}

func CloseWorkers() {
	threadPool.Close()
}

func Execute(f func() error) {
	threadPool.AddTask(f)
}
