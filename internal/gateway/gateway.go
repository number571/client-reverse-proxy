package gateway

import (
	"context"
	"io"
	"net"

	"client-reverse-proxy/internal/logger"
)

var (
	_ IGateway = &sGateway{}
)

type sGateway struct {
	fSettings ISettings
	fLogger   logger.ILogger
	fChConn   chan net.Conn
}

func NewGateway(pSettings ISettings, pLogger logger.ILogger) IGateway {
	return &sGateway{
		fSettings: pSettings,
		fLogger:   pLogger,
		fChConn:   make(chan net.Conn, pSettings.GetLimitClientProxyConns()),
	}
}

func (p *sGateway) GetSettings() ISettings {
	return p.fSettings
}

func (p *sGateway) GetLogger() logger.ILogger {
	return p.fLogger
}

func (p *sGateway) Run(pCtx context.Context) error {
	chErr := make(chan error, 1)

	go func() { chErr <- p.listenExternal(pCtx) }()
	go func() { chErr <- p.listenClientProxy(pCtx) }()

	select {
	case <-pCtx.Done():
		return pCtx.Err()
	case err := <-chErr:
		return err
	}
}

func (p *sGateway) listenExternal(pCtx context.Context) error {
	listener, err := net.Listen("tcp", p.fSettings.GetListenExternalAddr())
	if err != nil {
		p.fLogger.PushErro(err)
		return err
	}
	defer listener.Close()

	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		default:
			conn, err := listener.Accept()
			if err != nil {
				p.fLogger.PushWarn(err)
				continue
			}

			p.fLogger.PushInfo("accept connection")
			go p.handleConn(pCtx, conn)
		}
	}
}

func (p *sGateway) handleConn(pCtx context.Context, pConn net.Conn) {
	defer pConn.Close()

	var proxyConn net.Conn
	select {
	case conn := <-p.fChConn:
		p.fLogger.PushInfo("recv connection from channel")
		proxyConn = conn
	default:
		p.fLogger.PushWarn("empty connection channel")
		return
	}

	defer proxyConn.Close()

	go func() { _, _ = io.Copy(pConn, proxyConn) }()
	_, _ = io.Copy(proxyConn, pConn)
}

func (p *sGateway) listenClientProxy(ctx context.Context) error {
	listener, err := net.Listen("tcp", p.fSettings.GetListenClientProxyAddr())
	if err != nil {
		p.fLogger.PushErro(err)
		return err
	}
	defer listener.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			conn, err := listener.Accept()
			if err != nil {
				p.fLogger.PushWarn(err)
				continue
			}
			p.fLogger.PushInfo("send connection to channel")
			p.fChConn <- conn
		}
	}
}
