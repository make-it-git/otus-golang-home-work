package main

import (
	"context"
	"flag"
	"fmt"
	rabbit2 "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/sender"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	cfg "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/sender/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := cfg.NewSenderConfig(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rabbit, err := rabbit2.NewConsumer(config.Rabbit.Connection)
	if err != nil {
		logg.Error(err)
		os.Exit(1)
	}

	send := sender.New(logg, rabbit)
	go send.Run()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		rabbit.Stop(ctx)
		send.Stop()
	}()

	logg.Info("sender is running...")

	wg.Wait()
}
