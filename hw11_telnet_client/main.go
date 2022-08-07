package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func main() {
	log.SetOutput(os.Stderr)

	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Second*10, "Please provide timeout, default 10s")
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatal("Provide host and port")
	}
	host := flag.Args()[0]
	port, err := strconv.Atoi(flag.Args()[1])
	if err != nil {
		log.Fatal("Invalid port provided")
	}

	server := fmt.Sprintf("%s:%d", host, port)
	client := NewTelnetClient(server, timeout, os.Stdin, os.Stdout)
	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("...Connected to", server)

	ctx, cancel := context.WithCancel(context.Background())
	go gracefulShutdown(ctx, cancel, time.Second*3, syscall.SIGINT, syscall.SIGTERM)()

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go runUntilCompleted(ctx, cancel, wg, func() error {
		return client.Send()
	})
	go runUntilCompleted(ctx, cancel, wg, func() error {
		return client.Receive()
	})
	wg.Wait()
	_ = client.Close()
}
