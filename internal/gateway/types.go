package gateway

import (
	"client-reverse-proxy/internal/types"
	"time"

	"client-reverse-proxy/internal/logger"
)

type IGateway interface {
	types.IRunner

	GetSettings() ISettings
	GetLogger() logger.ILogger
}

type ISettings interface {
	GetDialTimeout() time.Duration
	GetListenExternalAddr() string
	GetListenClientProxyAddr() string
	GetLimitClientProxyConns() uint64
}
