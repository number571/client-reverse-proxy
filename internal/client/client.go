package client

import (
	"context"
	"io"
	"net"

	"client-reverse-proxy/internal/logger"
)

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fSettings ISettings
	fLogger   logger.ILogger
}

func NewClient(pSettings ISettings, pLogger logger.ILogger) IClient {
	return &sClient{
		fSettings: pSettings,
		fLogger:   pLogger,
	}
}

func (p *sClient) GetSettings() ISettings {
	return p.fSettings
}

func (p *sClient) GetLogger() logger.ILogger {
	return p.fLogger
}

func (p *sClient) Run(pCtx context.Context) error {
	connLimiter := make(chan struct{}, p.fSettings.GetNumGatewayConns())

	for {
		select {
		case <-pCtx.Done():
			return pCtx.Err()
		default:
			conn, err := p.getGatewayConn(pCtx)
			if err != nil {
				p.fLogger.PushWarn(err)
				<-connLimiter
				break // select
			}
			go func() {
				p.handleConn(pCtx, conn)
				p.fLogger.PushWarn("reconnecting")
				<-connLimiter
			}()
		}
	}
}

func (p *sClient) handleConn(pCtx context.Context, pConn net.Conn) {
	defer pConn.Close()

	d := net.Dialer{Timeout: p.fSettings.GetDialTimeout()}
	serverConn, err := d.DialContext(pCtx, "tcp", p.fSettings.GetServerAddr())
	if err != nil {
		p.fLogger.PushWarn(err)
		return
	}
	defer serverConn.Close()

	idleConn := newIdleTimeoutConn(pConn, p.fSettings.GetIdleTimeout())
	idleServerConn := newIdleTimeoutConn(serverConn, p.fSettings.GetIdleTimeout())

	chClosed := make(chan struct{}, 1)
	go func() {
		_, _ = io.Copy(idleConn, idleServerConn)
		chClosed <- struct{}{}
	}()
	go func() {
		_, _ = io.Copy(idleServerConn, idleConn)
		chClosed <- struct{}{}
	}()

	select {
	case <-pCtx.Done():
	case <-chClosed:
	}
}

func (p *sClient) getGatewayConn(ctx context.Context) (net.Conn, error) {
	d := net.Dialer{Timeout: p.fSettings.GetDialTimeout()}
	conn, err := d.DialContext(ctx, "tcp", p.fSettings.GetGatewayAddr())
	if err != nil {
		p.fLogger.PushWarn(err)
		return nil, err
	}
	p.fLogger.PushInfo("success connected to gateway")
	return conn, nil
}
