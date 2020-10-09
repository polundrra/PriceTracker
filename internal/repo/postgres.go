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

func (p *postgres) AddSubscription(ctx context.Context, adID uint64, email string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	tx, err := p.pool.BeginEx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error begin transaction, %v", err)
	}

	defer func(tx *pgx.Tx) {
		if err := tx.RollbackEx(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("error rollback AddSubscription, %v", err)
		}
	}(tx)

	var emailID uint64
	if err := tx.QueryRowEx(ctx, "insert into mail(email) values ($1) returning id", nil, email).Scan(&emailID); err != nil {
		return fmt.Errorf("error insert into mail, %v", err)
	}

	if _, err := tx.ExecEx(ctx, "insert into subscription(email_id, ad_id) values ($1, $2)", nil, emailID, adID); err != nil {
		return fmt.Errorf("err insert into subscription, %v", err)
	}

	if err := tx.CommitEx(ctx); err != nil {
		return fmt.Errorf("error commit AddSubscription transaction %v", err)
	}

	return nil
}

func (p *postgres) AddSubscriptionIfEmailExists(ctx context.Context, adID uint64, email string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var emailID uint64
	if err := p.pool.QueryRowEx(ctx, "select id from mail where email = $1", nil, email).Scan(&emailID); err != nil {
		return err
	}

	if _, err := p.pool.ExecEx(ctx, "insert into subscription(email_id, ad_id) values ($1, $2)", nil, emailID, adID); err != nil {
		return fmt.Errorf("err insert into subscription, %v", err)
	}

	return nil
}

func (p *postgres) GetEmailID(ctx context.Context, email string) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var id uint64
	if err := p.pool.QueryRowEx(ctx, "select id from mail where email = $1",nil,  email).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return 0, fmt.Errorf("error get email id, %v", err)
	}

	return id, nil
}

func (p *postgres) GetEmailsByAdID(ctx context.Context, adID uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var emails []string
	var email string

	rows, err := p.pool.QueryEx(ctx, "select email from mail m inner join subscription s on m.id = s.email_id " +
		"where ad_id = $1", nil, adID)
	if err != nil {
		return nil, fmt.Errorf("error get emails by ad id, %v", err)
	}

	for rows.Next() {
		if err := rows.Scan(&email); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("err scan rows in GetEmailByAdID, %v", err)
		}
		emails = append(emails, email)
	}

	return emails, nil
}



