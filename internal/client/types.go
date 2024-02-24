package client

import (
	"client-reverse-proxy/internal/types"
	"time"

	"client-reverse-proxy/internal/logger"
)

type IClient interface {
	types.IRunner

	GetSettings() ISettings
	GetLogger() logger.ILogger
}

type ISettings interface {
	GetServerAddr() string
	GetGatewayAddr() string
	GetIdleTimeout() time.Duration
	GetDialTimeout() time.Duration
	GetNumGatewayConns() uint64
}
