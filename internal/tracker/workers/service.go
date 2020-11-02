package workers

import (
	"context"
	"github.com/polundrra/PriceTracker/internal/tracker/priceinfo"
	"github.com/polundrra/PriceTracker/internal/tracker/repo"
	"time"
)

type Service interface {
	UpdatePrice(ctx context.Context, period time.Duration) error
	GetInfoForMailing(ctx context.Context, delay time.Duration)  ([]MessageInfo, error)
}

type MessageInfo struct {
	Emails []string
	Ad string
	NewPrice uint64
}

func New(repo repo.Repo, priceInfo priceinfo.Service) Service {
	return &workersService{
		repo: repo,
		priceClient: priceInfo,
	}
}