package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type errCount struct {
	mu    sync.Mutex
	count int
}

func (e *errCount) get() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.count
}

func (e *errCount) increment() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.count++
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workerCount, maxErrors int) error {
	if workerCount <= 0 {
		workerCount = 1
	}

	taskChan := make(chan Task)
	err := &errCount{}
	wg := &sync.WaitGroup{}

	startWorkers(workerCount, taskChan, err, wg)
	publishTasks(tasks, taskChan, err, wg, maxErrors)

	wg.Wait()
	if err.get() >= maxErrors {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func startWorkers(workerCount int, taskChan chan Task, err *errCount, wg *sync.WaitGroup) {
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(taskChan chan Task, err *errCount, wg *sync.WaitGroup) {
			defer wg.Done()
			for task := range taskChan {
				if e := task(); e != nil {
					err.increment()
				}
			}
		}(taskChan, err, wg)
	}
}

func publishTasks(tasks []Task, taskChan chan Task, err *errCount, wg *sync.WaitGroup, maxErrors int) {
	wg.Add(1)
	go func(tasks []Task, taskChan chan Task, err *errCount, wg *sync.WaitGroup, maxErrors int) {
		defer wg.Done()
		defer close(taskChan)

		for _, task := range tasks {
			taskChan <- task
			if err.get() >= maxErrors {
				break
			}
		}
	}(tasks, taskChan, err, wg, maxErrors)
}
