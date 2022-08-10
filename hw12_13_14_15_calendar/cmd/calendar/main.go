package main

import (
	"context"
	"flag"
	"fmt"
	cfg "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	storage2 "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
	sqlstorage "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/server/http"
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

	// storage := memorystorage.New()
	storage := sqlstorage.New(config.Storage.Connection)
	err = storage.Connect(context.Background())
	ev := storage2.Event{
		ID:               "A0EEBC99-9C0B-4EF8-BB6D-6BB9BD380A12",
		Title:            "test 2",
		StartTime:        time.Now(),
		Duration:         time.Hour * 3,
		Description:      nil,
		OwnerId:          33,
		NotificationTime: nil,
	}
	e := storage.Create(ev)
	fmt.Println(e)
	e = storage.Update(ev)
	fmt.Println(e)
	e = storage.Delete(ev.ID)
	fmt.Println(e)
	e = storage.Create(ev)
	fmt.Println(e)
	events, e := storage.ListMonth(ev.StartTime)
	fmt.Println(events, e)

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, &config.Http)

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
