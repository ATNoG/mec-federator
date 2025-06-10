package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mankings/mec-federator/internal/router"
)

// Task represents a periodic task that needs to be executed
type Task struct {
	Name     string
	Interval time.Duration
	Handler  func(ctx context.Context) error
}

// Scheduler manages periodic tasks
type Scheduler struct {
	tasks    []Task
	services *router.Services
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewScheduler creates a new scheduler instance
func NewScheduler(services *router.Services) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		services: services,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// AddTask adds a new periodic task to the scheduler
func (s *Scheduler) AddTask(task Task) {
	log.Printf("Adding task %s with interval %s", task.Name, task.Interval)
	s.tasks = append(s.tasks, task)
}

// Start begins executing all registered tasks
func (s *Scheduler) Start() {
	for _, task := range s.tasks {
		s.wg.Add(1)
		go func(t Task) {
			defer s.wg.Done()
			ticker := time.NewTicker(t.Interval)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := t.Handler(s.ctx); err != nil {
						log.Printf("Error executing task %s: %v", t.Name, err)
					}
				case <-s.ctx.Done():
					return
				}
			}
		}(task)
	}
}

// Stop gracefully stops all running tasks
func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
}
