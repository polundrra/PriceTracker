package service

import (
	"context"
	"github.com/polundrra/PriceTracker/internal/tracker/repo"
)

type Service interface {
	CreateSubscription(ctx context.Context, email string, ad uint64) error
}

func New(repo repo.Repo) Service {
	return &appService{
		repo: repo,
	}
}


