package connect

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAccount(ctx context.Context, acc *Account) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO accounts (user_id, type, country, email, business_type, status, platform_fee_percent) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		acc.UserID, acc.Type, acc.Country, acc.Email, acc.BusinessType, acc.Status, acc.PlatformFeePercent).
		Scan(&acc.ID, &acc.CreatedAt, &acc.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO payout_settings (account_id, interval, bank_account_id) 
		 VALUES ($1, $2, $3)`,
		acc.ID, acc.PayoutSettings.Interval, acc.PayoutSettings.BankAccountID)

	if err != nil {
		return fmt.Errorf("failed to create payout settings: %w", err)
	}

	return tx.Commit()
}

func (r *Repository) GetAccount(ctx context.Context, id string) (*Account, error) {
	var acc Account
	err := r.db.QueryRowContext(ctx,
		`SELECT a.id, a.user_id, a.type, a.country, a.email, a.business_type, a.status, a.platform_fee_percent, a.created_at, a.updated_at,
		        p.interval, p.bank_account_id
		 FROM accounts a
		 LEFT JOIN payout_settings p ON a.id = p.account_id
		 WHERE a.id = $1`,
		id).Scan(
		&acc.ID, &acc.UserID, &acc.Type, &acc.Country, &acc.Email, &acc.BusinessType, &acc.Status, &acc.PlatformFeePercent, &acc.CreatedAt, &acc.UpdatedAt,
		&acc.PayoutSettings.Interval, &acc.PayoutSettings.BankAccountID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	acc.PayoutSettings.AccountID = acc.ID
	return &acc, nil
}

func (r *Repository) UpdateAccount(ctx context.Context, acc *Account) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx,
		`UPDATE accounts SET email = $1, platform_fee_percent = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`,
		acc.Email, acc.PlatformFeePercent, acc.ID)
	if err != nil {
		return fmt.Errorf("failed to update account: %w", err)
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE payout_settings SET interval = $1, bank_account_id = $2, updated_at = CURRENT_TIMESTAMP WHERE account_id = $3`,
		acc.PayoutSettings.Interval, acc.PayoutSettings.BankAccountID, acc.ID)
	if err != nil {
		return fmt.Errorf("failed to update payout settings: %w", err)
	}

	return tx.Commit()
}
