package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func produce(ctx context.Context, chTasks chan<- Task, tasks []Task, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(chTasks)
	i := 0
	for i < len(tasks) {
		task := tasks[i]
		select {
		case <-ctx.Done():
			return
		case chTasks <- task:
			i++
		}
	}
}

func consume(ctx context.Context, errCh chan<- error, tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-tasks:
			if task == nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case errCh <- task():
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving maxErrors errors from tasks.
func Run(tasks []Task, n, maxErrors int) error {
	if maxErrors < 0 {
		return ErrErrorsLimitExceeded
	}

	runGoroutines := n
	if len(tasks) < runGoroutines {
		runGoroutines = len(tasks)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// it makes sense to set size of channel buffer no more than amount of goroutines
	chTasks := make(chan Task, runGoroutines)

	wgProducer := &sync.WaitGroup{}
	wgProducer.Add(1)
	go produce(ctx, chTasks, tasks, wgProducer)

	errCh := make(chan error)
	wgConsumer := &sync.WaitGroup{}
	wgConsumer.Add(runGoroutines)
	for i := 0; i < runGoroutines; i++ {
		go consume(ctx, errCh, chTasks, wgConsumer)
	}

	chResult := make(chan error)
	go handleResults(cancel, errCh, maxErrors, chResult)

	wgProducer.Wait()
	wgConsumer.Wait()
	close(errCh)

	resultError := <-chResult
	close(chResult)

	return resultError
}

func handleResults(cancel func(), errCh chan error, maxErrors int, chResult chan<- error) {
	totalErrors := 0

	for err := range errCh {
		if err != nil {
			totalErrors++
		}
		if totalErrors >= maxErrors {
			cancel()
			chResult <- ErrErrorsLimitExceeded
			return
		}
	}

	chResult <- nil
}
