package service

import (
	"context"
	"github.com/nordew/scope_test/internal/model"
)

type WorkerService interface {
	Submit(job model.Job)
	Shutdown()
	Wait()
	HandleErrors(ctx context.Context) error
	LogErrors(ctx context.Context)
}
