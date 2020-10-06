package service

import (
	"context"
	"github.com/polundrra/PriceTracker/internal/repo"
)

type appService struct {
	repo repo.Repo

}

func (s *appService) CreateSubscription(ctx context.Context, email, adURL string) error {
	emailID, err := s.repo.GetEmailID(ctx, email)
	if err != nil {
		return err
	}

	if emailID != 0 {
		adExists := false
		ads, err := s.repo.GetAdsByEmailID(ctx, emailID)
		if err != nil {
			return err
		}

		for i := range ads {
			if ads[i] == adURL {
				adExists = true
			}
		}

		if adExists {
			return ErrSubscriptionExists
		}

		if err := s.repo.AddSubscription(ctx, adURL, email, emailID); err != nil {
			return err
		}

		return nil
	}

	if err := s.repo.AddSubscription(ctx, adURL, email, emailID); err != nil {
		return err
	}

	return nil
}

