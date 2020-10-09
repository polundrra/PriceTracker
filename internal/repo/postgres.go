package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/polundrra/PriceTracker/internal/service"
	"log"
	"time"
)

type postgres struct {
	pool *pgx.ConnPool
	timeout time.Duration
}

func (p *postgres) GetInfoForMailing(ctx context.Context, dif time.Duration) ([]service.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	start := time.Now()
	t := start.Add(-dif)
	rows, err := p.pool.QueryEx(ctx, "select m.email, a.price from subscription s inner join mail m on m.id = s.mail_id" +
		"inner join advertisement a on s.ad_id = a.id where last_check_at >= $1", nil, t)
	if err != nil {
		return nil, fmt.Errorf("error get info, %v", err)
	}

	var messages []service.Message
	var message service.Message
	for rows.Next() {
		if err := rows.Scan(&message); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("err scan rows in GetInfoForMailing, %v", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (p *postgres) UpdatePrice(ctx context.Context, ad uint64, price string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if _, err := p.pool.ExecEx(ctx, "update advertisement set price = $1, last_check_at = $2 where ad = $3", nil, price, time.Now(), ad); err != nil {
		return fmt.Errorf("error update price, %v", err)
	}

	return nil
}

func (p *postgres) GetPriceByAd(ctx context.Context, ad uint64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var price string
	if err := p.pool.QueryRowEx(ctx, "select price from advertisement where ad = $1", nil, ad).Scan(&price); err != nil {
		return "", fmt.Errorf("error get price by ad id, %v", err)
	}

	return price, nil
}

func (p *postgres) AddSubscription(ctx context.Context, ad uint64, email string) error {
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
	if err := tx.QueryRowEx(ctx, "select id from mail where email = $1", nil, email).Scan(&emailID); err != nil {
		return fmt.Errorf("error select mail id, %v", err)
	}

	var adID uint64
	if err := tx.QueryRowEx(ctx, "select id from advertisement where ad = $1", nil, ad).Scan(&adID); err != nil {
		return fmt.Errorf("error select ad id, %v", err)
	}

	if _, err := tx.ExecEx(ctx, "insert into subscription(email_id, ad_id) values ($1, $2)", nil, emailID, adID); err != nil {
		return fmt.Errorf("err insert into subscription, %v", err)
	}

	if err := tx.CommitEx(ctx); err != nil {
		return fmt.Errorf("error commit AddSubscription transaction %v", err)
	}

	return nil
}

func (p *postgres) AddEmail(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if _, err := p.pool.ExecEx(ctx, "insert into mail(email) values ($1)",nil, email); err != nil {
		return fmt.Errorf("error insert into mail, %v", err)
	}

	return nil
}

func (p *postgres) AddAd(ctx context.Context, ad uint64, price string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if _, err := p.pool.ExecEx(ctx, "insert into advertisement(ad, price) values($1, $2)", nil, ad, price); err != nil {
		return fmt.Errorf("error insert into advertisement, %v", err)
	}

	return nil
}

func (p *postgres) GetAdID(ctx context.Context, ad uint64) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var id uint64
	if err := p.pool.QueryRowEx(ctx, "select id from advertisement where ad = $1",nil, ad).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}

		return 0, fmt.Errorf("error get ad id, %v", err)
	}

	return id, nil
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

func (p *postgres) GetEmailsByAd(ctx context.Context, ad uint64) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var emails []string
	var email string

	rows, err := p.pool.QueryEx(ctx, "select email from mail m inner join subscription s on m.id = s.email_id " +
		"where ad_id = $1", nil, ad)
	if err != nil {
		return nil, fmt.Errorf("error get emails by ad id, %v", err)
	}

	for rows.Next() {
		if err := rows.Scan(&email); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("err scan rows in GetEmailByAd, %v", err)
		}
		emails = append(emails, email)
	}

	return emails, nil
}


func (p *postgres) GetAdsForCheck(ctx context.Context, period time.Duration) ([]uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	now := time.Now()
	t := now.Add(-period)
	rows, err := p.pool.QueryEx(ctx, "select ad from advertisement where last_check_at <= $1", nil, t)
	if err != nil {
		return nil, fmt.Errorf("erorr get ads, %v", err)
	}

	var ads []uint64
	var ad uint64
	for rows.Next() {
		if err := rows.Scan(&ad); err != nil {
			if err == pgx.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("err scan rows in GetEmailByAd, %v", err)
		}
		ads = append(ads, ad)
	}

	return ads, nil
}




