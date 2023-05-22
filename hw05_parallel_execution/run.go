package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *SafeCounter) Get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// create safe counter
	safeCounter := SafeCounter{}

	// create channel
	jobs := make(chan Task, len(tasks)/2)

	// create wait group
	var wg sync.WaitGroup
	wg.Add(n)

	// create workers
	for i := 1; i <= n; i++ {
		go func() {
			defer wg.Done()
			for fu := range jobs {
				if safeCounter.Get() > m {
					return
				} else if fu() != nil {
					safeCounter.Inc()
				}
			}
		}()
	}

	// create jobs for workers
	for _, t := range tasks {
		jobs <- t
	}
	close(jobs)

	// wait workers
	wg.Wait()

	if safeCounter.value <= m+n {
		ErrErrorsLimitExceeded = nil
	}

	return ErrErrorsLimitExceeded
}
