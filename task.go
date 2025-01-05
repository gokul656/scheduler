package main

import (
	"context"
	"errors"
	"log"
	"time"
)

type Metrics interface {
	GetStartedAt() time.Time
	GetStoppedAt() time.Time
	Runtime() time.Duration
}

type Task func() error

type TaskMetrics struct {
	ctx        context.Context
	taskID     string
	enableLogs bool
	startedAt  time.Time
	stoppedAt  time.Time
	task       Task
}

func (tm *TaskMetrics) run() {
	tm.startedAt = time.Now()
	defer func() {
		tm.stoppedAt = time.Now()
		tm.logRuntime()
	}()

	tm.task()
}

func (tm *TaskMetrics) runWithTimeout(timeout time.Duration) error {
	tm.startedAt = time.Now()
	defer func() {
		tm.stoppedAt = time.Now()
		tm.logRuntime()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan error, 1)

	go func() {
		done <- tm.task()
	}()

	select {
	case <-ctx.Done():
		return errors.New("timeout: task took too long to respond")
	case err := <-done:
		return err
	}
}

func (tm *TaskMetrics) logRuntime() {
	if tm.enableLogs {
		log.Printf("runtime: %v\n", tm.Runtime())
	}
}

func NewMetrics(ctx context.Context, taskID string, task Task, enableLogs bool) *TaskMetrics {
	return &TaskMetrics{
		taskID:     taskID,
		ctx:        ctx,
		enableLogs: enableLogs,
		task:       task,
	}
}

func (tm TaskMetrics) GetStartedAt() time.Time { return tm.startedAt }
func (tm TaskMetrics) GetStoppedAt() time.Time { return tm.stoppedAt }
func (tm TaskMetrics) Runtime() time.Duration {
	return tm.GetStoppedAt().Sub(tm.GetStartedAt())
}
