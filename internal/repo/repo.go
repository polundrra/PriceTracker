package repo

import (
	"context"
	"github.com/jackc/pgx"
	"time"
)

type Repo interface {
	AddSubscription(ctx context.Context, adURL string, email string, emailID uint64) error
	GetEmailID(ctx context.Context, email string) (uint64, error)
	GetAdsByEmailID(ctx context.Context, emailID uint64) ([]string, error)
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