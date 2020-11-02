package priceinfo

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

type client struct {
	client  fasthttp.Client
	addr    string
	key     string
	timeout time.Duration
}

func (c *client) GetPriceByAd(ad uint64) (uint64, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(c.addr + strconv.FormatUint(ad, 10) + c.key)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.client.DoTimeout(req, resp, c.timeout); err != nil {
		return 0, fmt.Errorf("Client get failed: %s\n", err)
	}

	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("Expected status code %d, but got %d\n", fasthttp.StatusOK, resp.StatusCode())
	}

	body := resp.Body()
	value, err := jsonparser.GetString(body, "price", "value")
	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}

	price *= 100 //перевод в копейки

	return price, nil
}
