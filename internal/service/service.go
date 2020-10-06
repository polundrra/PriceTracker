package service

import (
	"context"
	"github.com/polundrra/PriceTracker/internal/repo"
)

type Service interface {
	CreateSubscription(ctx context.Context, email, adURL string) error
}

func New(repo repo.Repo) Service {
	return &appService{
		repo: repo,
	}
}
