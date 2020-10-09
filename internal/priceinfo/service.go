package priceinfo

import (
	"github.com/valyala/fasthttp"
	"time"
)

type clientService interface {
	GetPriceByAdID(adID uint64) (uint64, error)
}

type Opts struct {
	Addr string
	Timeout int
}

func New(opts Opts) clientService {
	return &client{
		client: fasthttp.Client {},
		addr: opts.Addr,
		timeout: time.Duration(opts.Timeout) * time.Second,
	}
}