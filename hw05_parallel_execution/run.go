package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workerCount, maxErrors int) error {
	taskChan := make(chan Task, workerCount)
	resultChan := make(chan error, workerCount)
	shutdownChan := make(chan struct{})

	defer func(shutdownChan chan struct{}) {
		for i := 0; i < workerCount; i++ {
			shutdownChan <- struct{}{}
		}
		close(shutdownChan)
		close(taskChan)
		close(resultChan)
	}(shutdownChan)

	startWorkers(workerCount, taskChan, resultChan, shutdownChan)

	return runTasksAndReceiveResults(tasks, taskChan, resultChan, maxErrors)
}

func startWorkers(n int, taskChan chan Task, resultChan chan error, shutdownChan chan struct{}) {
	for i := 0; i < n; i++ {
		go func(taskChan chan Task, resultChan chan error, shutdownChan chan struct{}) {
			for {
				select {
				case task := <-taskChan:
					select {
					case resultChan <- task():
						// result published
					case <-shutdownChan:
						return
					}
				case <-shutdownChan:
					return
				}
			}
		}(taskChan, resultChan, shutdownChan)
	}
}

func runTasksAndReceiveResults(
	tasks []Task,
	taskChan chan Task,
	resultChan chan error,
	maxErrors int,
) error {
	numStarted, numResults, numErrors := 0, 0, 0
	for {
		tasksToStart := cap(taskChan) - len(taskChan)
		for i, pos := 0, numStarted; i < tasksToStart && pos+i < len(tasks); i++ {
			taskChan <- tasks[pos+i]
			numStarted++
		}

		for {
			result := <-resultChan
			if result != nil {
				numErrors++
			}
			numResults++
			// drain channel to correctly check for max errors
			if numResults == numStarted {
				break
			}
		}

		if numErrors >= maxErrors {
			return ErrErrorsLimitExceeded
		}
		if numResults >= len(tasks) {
			return nil
		}
	}
}
