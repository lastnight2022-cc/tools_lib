package worker_pool

import (
	"errors"
	"sync"
)

// Task represents a function that can be submitted to the worker pool.
type Task func() error

// WorkerPool manages a pool of workers that execute tasks concurrently.
type WorkerPool struct {
	taskQueue chan Task
	wg        sync.WaitGroup
	closeChan chan struct{}
}

// NewWorkerPool creates a new WorkerPool with the specified number of workers.
func NewWorkerPool(workerCount int) (*WorkerPool, error) {
	if workerCount <= 0 {
		return nil, errors.New("worker count must be positive")
	}

	wp := &WorkerPool{
		taskQueue: make(chan Task),
		closeChan: make(chan struct{}),
	}
	wp.startWorkers(workerCount)
	return wp, nil
}

// startWorkers starts the specified number of workers.
func (wp *WorkerPool) startWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker()
	}
}

// worker is the function that each worker goroutine executes.
func (wp *WorkerPool) worker() {
	defer wp.wg.Done()
	for {
		select {
		case task, ok := <-wp.taskQueue:
			if !ok {
				// The channel has been closed. Exit the worker.
				return
			}
			if err := task(); err != nil {
				// Handle the error as appropriate for your application.
				// Here we simply print it out.
				println("Task failed:", err.Error())
			}
		case <-wp.closeChan:
			return
		}
	}
}

// Submit adds a new task to the queue.
func (wp *WorkerPool) Submit(task Task) error {
	select {
	case wp.taskQueue <- task:
		return nil
	case <-wp.closeChan:
		return errors.New("worker pool is closed")
	}
}

// Close stops all workers and waits for them to finish processing their current tasks.
func (wp *WorkerPool) Close() {
	close(wp.closeChan)
	close(wp.taskQueue)
	wp.wg.Wait()
}
