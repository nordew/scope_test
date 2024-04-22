package service

import (
	"context"
	"errors"
	"github.com/nordew/scope_test/internal/model"
	"log"
	"sync"
)

type WorkerPool struct {
	workers    []*Worker
	jobQueue   chan chan model.Job
	maxWorkers int
	wg         sync.WaitGroup
	errors     chan error
	mu         sync.Mutex
}

func NewWorkerPool(maxWorkers, maxQueueSize int) *WorkerPool {
	jobQueue := make(chan chan model.Job, maxQueueSize)
	errors := make(chan error)

	pool := &WorkerPool{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		errors:     errors,
	}

	for i := 0; i < maxWorkers; i++ {
		worker := NewWorker(jobQueue)

		pool.mu.Lock()
		pool.workers = append(pool.workers, worker)
		pool.mu.Unlock()

		go worker.Start(context.Background())
	}

	return pool
}

func (p *WorkerPool) Submit(job model.Job) {
	select {
	case jobChan := <-p.jobQueue:
		jobChan <- job
	default:
		p.errors <- errors.New("job queue is full")
	}
}

func (p *WorkerPool) Shutdown() {
	close(p.jobQueue)
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

func (p *WorkerPool) HandleErrors(ctx context.Context) error {
	var wg sync.WaitGroup

	for _, worker := range p.workers {
		wg.Add(1)
		go func(w *Worker) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case <-w.quit:
				return
			}
		}(worker)
	}

	wg.Wait()
	return nil
}

func (p *WorkerPool) LogErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-p.errors:
			log.Println("Error:", err)
		}
	}
}
