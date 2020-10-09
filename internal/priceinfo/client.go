package priceinfo

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

type client struct {
	client  fasthttp.Client
	addr    string
	key     string
	timeout time.Duration
}

func (c *client) GetPriceByAd(ad uint64) (string, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(c.addr + string(ad) + c.key)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.client.DoTimeout(req, resp, c.timeout); err != nil {
		log.Println(err, c.addr, c.timeout)
		return "", fmt.Errorf("Client get failed: %s\n", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("Expected status code %d but got %d\n", fasthttp.StatusOK, resp.StatusCode())
	}

	body := resp.Body()
	price, err := jsonparser.GetString(body, "price", "value")
	if err != nil {
		return "", err
	}

	return price, nil
}
