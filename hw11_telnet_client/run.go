package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func gracefulShutdown(
	ctx context.Context, cancel context.CancelFunc, timeout time.Duration, signals ...os.Signal) func() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	return func() {
		select {
		case s := <-c:
			log.Println("...stop via signal: ", s)
			cancel()
		case <-ctx.Done():
			log.Println("...stop via context")
			break
		}
		time.Sleep(timeout)
		os.Exit(0)
	}
}

func runUntilCompleted(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, f func() error) {
	ok := true
	for ok {
		select {
		case <-ctx.Done():
			ok = false
		default:
			err := f()
			if errors.Is(err, ErrCompleted) {
				ok = false
			} else if err != nil {
				log.Println("...error received", err)
				ok = false
			}
		}
	}
	cancel()
	wg.Done()
}
