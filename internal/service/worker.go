package service

import (
	"context"
	"github.com/nordew/scope_test/internal/model"
	"log"
)

type Worker struct {
	workerPool chan chan model.Job
	jobChannel chan model.Job
	quit       chan bool
	errors     chan error
}

func NewWorker(workerPool chan chan model.Job) *Worker {
	return &Worker{
		workerPool: workerPool,
		jobChannel: make(chan model.Job),
		quit:       make(chan bool),
		errors:     make(chan error),
	}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				w.workerPool <- w.jobChannel

				select {
				case job := <-w.jobChannel:
					err := job.Process()
					if err != nil {
						w.errors <- err
					}
				case <-w.quit:
					return
				}
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}

func (w *Worker) ErrorChannel() <-chan error {
	return w.errors
}

func (w *Worker) HandleErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-w.errors:
			log.Println("Worker error:", err)
		}
	}
}
