package service

import (
	"context"
	"github.com/polundrra/PriceTracker/internal/repo"
)

type appService struct {
	repo repo.Repo

}

func (s *appService) CreateSubscription(ctx context.Context, email string, adID uint64) error {
	emailID, err := s.repo.GetEmailID(ctx, email)
	if err != nil {
		return err
	}

	if emailID != 0 {
		adExists := false
		emails, err := s.repo.GetEmailsByAdID(ctx, adID)
		if err != nil {
			return err
		}

		for i := range emails {
			if emails[i] == email {
				adExists = true
			}
		}

		if adExists {
			return ErrSubscriptionExists
		}

		if err := s.repo.AddSubscriptionIfEmailExists(ctx, adID, email); err != nil {
			return err
		}


		return nil
	}

	if err := s.repo.AddSubscription(ctx, adID, email); err != nil {
		return err
	}

	return nil
}

