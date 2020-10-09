package repo

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/polundrra/PriceTracker/internal/service"
	"time"
)

type Repo interface {
	AddSubscription(ctx context.Context, ad uint64, email string) error
	AddEmail(ctx context.Context, email string) error
	AddAd(ctx context.Context, ad uint64, price string) error
	GetEmailID(ctx context.Context, email string) (uint64, error)
	GetAdID(ctx context.Context, ad uint64) (uint64, error)
	GetEmailsByAd(ctx context.Context, ad uint64) ([]string, error)
	GetPriceByAd(ctx context.Context, ad uint64) (string, error)
	UpdatePrice(ctx context.Context, ad uint64, price string) error
	GetAdsForCheck(ctx context.Context, period time.Duration) ([]uint64, error)
	GetInfoForMailing(ctx context.Context, dif time.Duration) ([]service.Message, error)
}

type Opts struct {
	Host string
	Port uint16
	Database string
	User string
	Password string
	Timeout int
}

func New(opts Opts) (Repo, error) {
	ConnConfig := pgx.ConnConfig{
		Host: opts.Host,
		Port: opts.Port,
		Database: opts.Database,
		User: opts.User,
		Password: opts.Password,
	}
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: ConnConfig,
	})
	if err != nil {
		return nil, err
	}
	repo := postgres{
		pool: pool,
		timeout: time.Duration(opts.Timeout) * time.Second,
	}
	return &repo, nil
}
