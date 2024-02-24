package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"client-reverse-proxy/internal/client"
	"client-reverse-proxy/internal/logger"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	closed := make(chan struct{})
	defer func() {
		cancel()
		<-closed
	}()

	go func() {
		defer func() { closed <- struct{}{} }()
		client := newClient()
		if err := client.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Fatal(err)
		}
	}()

	fmt.Println("client is listening...")
	defer fmt.Println("client is stopped...")

	<-shutdown
}

func newClient() client.IClient {
	return client.NewClient(
		client.NewSettings(&client.SSettings{
			FServerAddr:      "server:8080",
			FGatewayAddr:     "gateway:9001",
			FIdleTimeout:     time.Minute,
			FDialTimeout:     5 * time.Second,
			FNumGatewayConns: 1,
		}),
		logger.NewLogger(
			logger.NewSettings(&logger.SSettings{
				// FInfo: os.Stdout,
				FWarn: os.Stderr,
				FErro: os.Stderr,
			}),
			func(ia logger.ILogArg) string {
				switch x := ia.(type) {
				case string:
					return x
				case error:
					return x.Error()
				default:
					panic("unknown log type")
				}
			},
		),
	)
}
