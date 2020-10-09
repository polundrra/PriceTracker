package cmd

import (
	"github.com/polundrra/PriceTracker/internal/priceinfo"
	"github.com/polundrra/PriceTracker/internal/repo"
)

type conf struct {
	ServerPort string
	ReadTimeout int
	IdleTimeout int
	WriteTimeout int
	RepoOpts repo.Opts
	ClientOpts priceinfo.Opts
}
