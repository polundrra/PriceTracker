package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"time"
)

type postgres struct {
	pool *pgx.ConnPool
	timeout time.Duration
}

func (p *postgres) AddSubscription(ctx context.Context, adURL string, email string, emailID uint64) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if emailID != 0 {
		if _, err := p.pool.ExecEx(ctx, "insert into subscription(email_id, advertisement) values ($1, $2)", nil, emailID, adURL); err != nil {
			return fmt.Errorf("error insert adt with existing email, %v", err)
		}

		return nil
	}

	tx, err := p.pool.BeginEx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error begin transaction, %v", err)
	}

	defer func(tx *pgx.Tx) {
		if err := tx.RollbackEx(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("error rollback AddSubscription, %v", err)
		}
	}(tx)

	var id uint64
	if err := tx.QueryRowEx(ctx,"insert into email(email) values ($1) returning id", nil, email).Scan(&id); err != nil {
		return fmt.Errorf("err insert into email, %v", err)
	}

	if _, err := tx.ExecEx(ctx, "insert into subscription(email_id, advertisement) values ($1, $2)", nil, id, adURL); err != nil {
		return fmt.Errorf("err insert adt %s, %v", adURL, err)
	}

	if err := tx.CommitEx(ctx); err != nil {
		return fmt.Errorf("error commit AddSubscription transaction %v", err)
	}

	return nil
}

func (p *postgres) GetEmailID(ctx context.Context, email string) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var id uint64
	if err := p.pool.QueryRowEx(ctx, "select id from email where email = $1",nil,  email).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return 0, fmt.Errorf("error get email id, %v", err)
	}

	return id, nil
}

func (p *postgres) GetAdsByEmailID(ctx context.Context, emailID uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var ads []string
	var ad string

	rows, err := p.pool.QueryEx(ctx, "select advertisement from subscription where email_id = $1",nil,  emailID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("error get ads by email id, %v", err)
	}

	for rows.Next() {
		if err := rows.Scan(&ad); err != nil {
			return nil, fmt.Errorf("error scan rows in GetAdsByEmailID, %v", err)
		}
		ads = append(ads, ad)
	}

	return ads, nil
}



