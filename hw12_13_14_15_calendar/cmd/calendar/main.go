package main

import (
	"context"
	"flag"
	"fmt"
	cfg "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := cfg.NewConfig(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	storage := memorystorage.New()
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
