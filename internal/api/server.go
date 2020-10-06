package api

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/polundrra/PriceTracker/internal/service"
	"github.com/valyala/fasthttp"
	"net"
	"regexp"
	"strings"
)


type Server struct {
	service service.Service
}

func New(service service.Service) Server {
	return Server{service: service}
}

func (s *Server) Router() fasthttp.RequestHandler {
	r := router.New()
	r.GET("/subscription", s.subscribe)
	return r.Handler
}

type request struct {
	adURL string `json:"ad"`
	email string `json:"email"`
}

func (s *Server) subscribe(ctx *fasthttp.RequestCtx) {
	var req request
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Response.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.Write([]byte("error unmarshal request body" + err.Error()))
		return
	}

	//ВАЛИДАЦИЯ ССЫЛКИ ОБЪЯВЛЕНИЯ??

	if !isValidEmail(req.email) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString("invalid email")
		return
	}

	if err := s.service.CreateSubscription(ctx, req.email, req.adURL); err != nil {
		if err == service.ErrSubscriptionExists {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Write([]byte(err.Error()))
		return
	}

	ctx.Response.SetStatusCode(fasthttp.StatusOK)
}

func isValidEmail(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(e) < 3 && len(e) > 254 {
		return false
	}

	if !emailRegex.MatchString(e) {
		return false
	}

	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}

	return true
}
