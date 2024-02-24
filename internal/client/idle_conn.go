package client

import (
	"io"
	"net"
	"time"
)

type sIdleTimeoutConn struct {
	fConn    net.Conn
	fTimeout time.Duration
}

func newIdleTimeoutConn(pConn net.Conn, pTimeout time.Duration) io.ReadWriter {
	return &sIdleTimeoutConn{
		fConn:    pConn,
		fTimeout: pTimeout,
	}
}

func (p *sIdleTimeoutConn) Read(buf []byte) (int, error) {
	p.fConn.SetDeadline(time.Now().Add(p.fTimeout))
	return p.fConn.Read(buf)
}

func (p *sIdleTimeoutConn) Write(buf []byte) (int, error) {
	p.fConn.SetDeadline(time.Now().Add(p.fTimeout))
	return p.fConn.Write(buf)
}
