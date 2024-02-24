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

	"client-reverse-proxy/internal/gateway"
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
		gateway := newGateway()
		if err := gateway.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Fatal(err)
		}
	}()

	fmt.Println("gateway is listening...")
	defer fmt.Println("gateway is stopped...")

	<-shutdown
}

func newGateway() gateway.IGateway {
	return gateway.NewGateway(
		gateway.NewSettings(&gateway.SSettings{
			FDialTimeout:           5 * time.Second,
			FListenExternalAddr:    "0.0.0.0:9000",
			FListenClientProxyAddr: "0.0.0.0:9001",
			FLimitClientProxyConns: 128,
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
