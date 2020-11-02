package main

import (
	"github.com/polundrra/PriceTracker/internal/tracker/priceinfo"
	"github.com/polundrra/PriceTracker/internal/tracker/repo"
)

type conf struct {
	ServerPort   string
	ReadTimeout  int
	IdleTimeout  int
	WriteTimeout int
	RepoOpts     repo.Opts
	ClientOpts   priceinfo.Opts
}
