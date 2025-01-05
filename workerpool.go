package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type WorkerPool struct {
	ctx       context.Context
	wg        *sync.WaitGroup
	workers   int
	taskChan  chan Task
	startedAt time.Time
	stoppedAt time.Time
}

func (wp *WorkerPool) GetMaxWorkers() int {
	return wp.workers
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	return &WorkerPool{
		ctx:      context.Background(),
		workers:  maxWorkers,
		wg:       &sync.WaitGroup{},
		taskChan: make(chan Task, maxWorkers),
	}
}

func (wp *WorkerPool) handleTask(_ string, wg *sync.WaitGroup, taskChan chan Task) {
	defer wg.Done()

	time.Sleep(time.Second)

	for task := range taskChan {
		metrics := NewMetrics(wp.ctx, uuid.NewString(), task, true)
		if err := metrics.runWithTimeout(time.Duration(time.Second * 2)); err != nil {
			log.Printf("error %v", err)
			return
		}
	}
}

func (wp *WorkerPool) Spawn() {
	wg := wp.wg
	wg.Add(wp.GetMaxWorkers())

	for i := range wp.GetMaxWorkers() {
		name := fmt.Sprintf("worker_%d", i)
		go wp.handleTask(name, wp.wg, wp.taskChan)
	}

	wp.startedAt = time.Now()
}

func (wp *WorkerPool) Submit(task Task) {
	wp.taskChan <- task
}

func (wp *WorkerPool) Stop() {
	close(wp.taskChan)
	wp.wg.Wait()

	wp.stoppedAt = time.Now()
	log.Printf("workerpool runtime %v\n", wp.Runtime())

}

func (tm WorkerPool) GetStartedAt() time.Time { return tm.startedAt }
func (tm WorkerPool) GetStoppedAt() time.Time { return tm.stoppedAt }
func (tm WorkerPool) Runtime() time.Duration {
	return tm.GetStoppedAt().Sub(tm.GetStartedAt())
}
