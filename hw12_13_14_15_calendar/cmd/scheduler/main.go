package main

import (
	"context"
	"flag"
	"fmt"
	rabbit2 "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/rabbit"
	scheduler2 "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/scheduler"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cfg "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	sqlstorage "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/scheduler/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := cfg.NewSchedulerConfig(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rabbit, err := rabbit2.NewProducer(config.Rabbit.Connection)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	var storage scheduler2.Storage
	s := sqlstorage.New(config.Storage.Connection)
	err = s.Connect(context.Background())
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}
	storage = s
	var storageCloser interface {
		Close(ctx context.Context) error
	}
	storageCloser = s

	sched := scheduler2.New(logg, storage, rabbit, config.Rabbit.Timer, config.Cleanup)
	go sched.Run()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer wg.Done()
		defer storageCloser.Close(ctx)

		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		rabbit.Stop(ctx)
		sched.Stop()
	}()

	logg.Info("scheduler is running...")

	wg.Wait()
}
