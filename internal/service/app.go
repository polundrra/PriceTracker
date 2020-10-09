package service

import (
	"context"
	"github.com/polundrra/PriceTracker/internal/priceinfo"
	"github.com/polundrra/PriceTracker/internal/repo"
	"time"
)

type appService struct {
	repo repo.Repo
	priceClient priceinfo.Service
}


type Message struct {
	email string
	newPrice string
}

func (s *appService) GetInfoForMailing(ctx context.Context) ([]Message, error) {
	t := time.Second * 10
	info, err := s.repo.GetInfoForMailing(ctx, t)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}

	return info, nil
}

func (s *appService) UpdatePrice(ctx context.Context) error {
	t := time.Hour
	ads, err := s.repo.GetAdsForCheck(ctx, t)
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



func (s *appService) CreateSubscription(ctx context.Context, email string, ad uint64) error {
	emails, err := s.repo.GetEmailsByAd(ctx, ad)
	if err != nil {
		return err
	}

	for i := range emails {
		if emails[i] == email {
			return ErrSubscriptionExists
		}
	}

	emailID, err := s.repo.GetEmailID(ctx, email)
	if err != nil {
		return err
	}

	if emailID == 0 {
		if err := s.repo.AddEmail(ctx, email); err != nil {
			return err
		}
	}

	adID, err := s.repo.GetAdID(ctx, ad)
	if err != nil {
		return err
	}

	price, err := s.priceClient.GetPriceByAd(ad)
	if err != nil {
		return err
	}

	if adID == 0 {
		if err := s.repo.AddAd(ctx, ad, price); err != nil {
			return err
		}
	}

	repoPrice, err := s.repo.GetPriceByAd(ctx, ad)
	if err != nil {
		return err
	}

	if price != repoPrice {
		//делаем updatePrice
		if err := s.repo.UpdatePrice(ctx, ad, price); err != nil {
			return err
		}
		//тут рассылка
	}

	if err := s.repo.AddSubscription(ctx, ad, email); err != nil {
		return err
	}

	return nil
}


