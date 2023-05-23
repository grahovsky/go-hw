package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task, len(tasks)/2)

	if m < 0 {
		m = 0
	}
	maxErr := make(chan struct{}, m)
	for i := 1; i <= m; i++ {
		maxErr <- struct{}{}
	}
	close(maxErr)

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 1; i <= n; i++ {
		go func() {
			defer wg.Done()
			for fu := range jobs {
				if fu() != nil {
					if _, ok := <-maxErr; !ok {
						return
					}
				}
			}
		}()
	}

	for _, t := range tasks {
		if len(maxErr) > 0 {
			jobs <- t
		}
	}
	close(jobs)
	wg.Wait()

	if _, ok := <-maxErr; ok {
		ErrErrorsLimitExceeded = nil
	}

	return ErrErrorsLimitExceeded
}
