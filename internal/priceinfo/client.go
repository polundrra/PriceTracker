package priceinfo

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"time"
)

type client struct {
	client   fasthttp.Client
	addr     string
	timeout  time.Duration
}

func (c *client) GetPriceByAdID(adID uint64) (uint64, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(c.addr + string(adID))
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := c.client.DoTimeout(req, resp, c.timeout); err != nil {
		log.Println(err, c.addr, c.timeout)
		return 0, fmt.Errorf("Client get failed: %s\n", err)
	}

	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("Expected status code %d but got %d\n", fasthttp.StatusOK, resp.StatusCode())
	}

	body := resp.Body()
	value, err := jsonparser.GetString(body, "price", "value")
	if err != nil {
		return 0, err
	}

	price, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return uint64(price), nil
}
