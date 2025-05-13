package util

import (
	"sync"
)

type WorkerPool struct {
	taskChan chan func() error
	wg       sync.WaitGroup
	once     sync.Once
}

var threadPool = GetWorkerPool(4)

// GetWorkerPool 获取全局线程池实例(单例模式)
func GetWorkerPool(numWorkers int) *WorkerPool {
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
				LOGGER.Error("Function task failed: " + err.Error())
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
