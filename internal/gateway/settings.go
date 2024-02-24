package gateway

import "time"

var (
	_ ISettings = &sSettings{}
)

type SSettings sSettings
type sSettings struct {
	FDialTimeout           time.Duration
	FListenExternalAddr    string
	FListenClientProxyAddr string
	FLimitClientProxyConns uint64
}

func NewSettings(pSett *SSettings) ISettings {
	return (&sSettings{
		FDialTimeout:           pSett.FDialTimeout,
		FListenExternalAddr:    pSett.FListenExternalAddr,
		FListenClientProxyAddr: pSett.FListenClientProxyAddr,
		FLimitClientProxyConns: pSett.FLimitClientProxyConns,
	}).mustNotNull()
}

func (p *sSettings) mustNotNull() ISettings {
	if p.FDialTimeout == 0 {
		panic(`p.FDialTimeout == 0`)
	}
	if p.FListenExternalAddr == "" {
		panic(`p.FListenExternalAddr == ""`)
	}
	if p.FListenClientProxyAddr == "" {
		panic(`p.FListenClientProxyAddr == ""`)
	}
	if p.FLimitClientProxyConns == 0 {
		panic(`p.FLimitClientProxyConns == 0`)
	}
	return p
}

func (p *sSettings) GetDialTimeout() time.Duration {
	return p.FDialTimeout
}

func (p *sSettings) GetListenExternalAddr() string {
	return p.FListenExternalAddr
}

func (p *sSettings) GetListenClientProxyAddr() string {
	return p.FListenClientProxyAddr
}

func (p *sSettings) GetLimitClientProxyConns() uint64 {
	return p.FLimitClientProxyConns
}
