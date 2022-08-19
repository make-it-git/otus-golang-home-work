package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/server/grpc"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/make-it-git/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage/sql"
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

	var storage app.Storage
	var storageCloser interface {
		Close(ctx context.Context) error
	}
	switch config.Storage.Kind {
	case "db":
		s := sqlstorage.New(config.Storage.Connection)
		err = s.Connect(context.Background())
		if err != nil {
			logg.Error(err)
			os.Exit(1)
		}
		storage = s
		storageCloser = s
	case "memory":
		s := memorystorage.New()
		storage = s
	}

	calendar := app.New(logg, storage)
	srvHTTP := internalhttp.NewServer(logg, &config.HTTP, calendar)
	srvGRPC := grpc.NewServer(logg, &config.GRPC, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := new(sync.WaitGroup)
	wg.Add(3)

	go func() {
		defer wg.Done()

		if storageCloser != nil {
			defer func() {
				err := storageCloser.Close(ctx)
				if err != nil {
					logg.Error("failed to close storage connection: " + err.Error())
				}
			}()
		}

		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := srvHTTP.Stop(ctx); err != nil {
			logg.Error("failed to stop http srvHTTP: " + err.Error())
		}
		if err := srvGRPC.Stop(ctx); err != nil {
			logg.Error("failed to stop grpc srvGRPC: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		defer wg.Done()
		if err := srvHTTP.Start(ctx); err != nil {
			logg.Error("failed to start http srvHTTP: " + err.Error())
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		if err := srvGRPC.Start(ctx); err != nil {
			logg.Error("failed to start http srvGRPC: " + err.Error())
			cancel()
		}
	}()
	wg.Wait()
}
