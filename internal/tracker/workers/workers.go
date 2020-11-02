package workers

import (
	"context"
	"github.com/polundrra/PriceTracker/internal/tracker/priceinfo"
	"github.com/polundrra/PriceTracker/internal/tracker/repo"
	"time"
)

type workersService struct {
	repo        repo.Repo
	priceClient priceinfo.Service
}

func (s *workersService) GetInfoForMailing(ctx context.Context, delay time.Duration) ([]MessageInfo, error) {
	panic("implement me")
}

func (s *workersService) UpdatePrice(ctx context.Context, period time.Duration) error {
	ads, err := s.repo.GetAdsByLastCheck(ctx, period)
	if err != nil {
		return err
	}

	for i := range ads {
		price, err := s.priceClient.GetPriceByAd(ads[i])
		if err != nil {
			return err
		}

		repoPrice, err := s.repo.GetPriceByAd(ctx, ads[i]);
		if err != nil {
			return nil
		}

		if price != repoPrice {
			if err := s.repo.UpdatePrice(ctx, ads[i], price); err != nil {
				return err
			}
		}
	}
	return nil
}
