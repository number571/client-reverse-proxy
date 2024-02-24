package client

import "time"

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FServerAddr      string
	FGatewayAddr     string
	FIdleTimeout     time.Duration
	FDialTimeout     time.Duration
	FNumGatewayConns uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FServerAddr:      pSett.FServerAddr,
		FGatewayAddr:     pSett.FGatewayAddr,
		FIdleTimeout:     pSett.FIdleTimeout,
		FDialTimeout:     pSett.FDialTimeout,
		FNumGatewayConns: pSett.FNumGatewayConns,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FServerAddr == "" {
		panic(`p.FServerAddr == ""`)
	}
	if p.FGatewayAddr == "" {
		panic(`p.FGatewayAddr == ""`)
	}
	if p.FIdleTimeout == 0 {
		panic(`p.FIdleTimeout == 0`)
	}
	if p.FDialTimeout == 0 {
		panic(`p.FDialTimeout == 0`)
	}
	if p.FNumGatewayConns == 0 {
		panic(`p.FNumGatewayConns == 0`)
	}
	return p
}

func (p *sSettings) GetServerAddr() string {
	return p.FServerAddr
}

func (p *sSettings) GetGatewayAddr() string {
	return p.FGatewayAddr
}

func (p *sSettings) GetIdleTimeout() time.Duration {
	return p.FIdleTimeout
}

func (p *sSettings) GetDialTimeout() time.Duration {
	return p.FDialTimeout
}

func (p *sSettings) GetNumGatewayConns() uint64 {
	return p.FNumGatewayConns
}
