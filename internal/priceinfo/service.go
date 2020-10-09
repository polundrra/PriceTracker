package priceinfo

import (
	"github.com/valyala/fasthttp"
	"time"
)

type Service interface {
	GetPriceByAd(ad uint64) (string, error)
}

type Opts struct {
	Addr string
	Timeout int
}

func New(opts Opts) Service {
	return &client{
		client: fasthttp.Client {},
		addr: opts.Addr,
		timeout: time.Duration(opts.Timeout) * time.Second,
	}
}