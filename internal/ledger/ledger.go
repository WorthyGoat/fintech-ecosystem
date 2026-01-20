package ledger

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type AccountType string
type TransactionType string

const (
	Asset     AccountType = "asset"
	Liability AccountType = "liability"
	Equity    AccountType = "equity"
	Revenue   AccountType = "revenue"
	Expense   AccountType = "expense"
)

const (
	Debit  TransactionType = "debit"
	Credit TransactionType = "credit"
)

type Account struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Type      AccountType `json:"type"`
	Balance   int64       `json:"balance"`
	UserID    *string     `json:"user_id,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
}

type Transaction struct {
	ID          string    `json:"id"`
	ReferenceID string    `json:"reference_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Entry struct {
	ID            string          `json:"id"`
	TransactionID string          `json:"transaction_id"`
	AccountID     string          `json:"account_id"`
	Amount        int64           `json:"amount"`
	Direction     TransactionType `json:"direction"`
	CreatedAt     time.Time       `json:"created_at"`
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateAccount(ctx context.Context, name string, accType AccountType, userID *string) (*Account, error) {
	acc := &Account{
		Name:   name,
		Type:   accType,
		UserID: userID,
	}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO accounts (name, type, user_id) VALUES ($1, $2, $3) RETURNING id, balance, created_at`,
		name, accType, userID).Scan(&acc.ID, &acc.Balance, &acc.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}
	return acc, nil
}

func (r *Repository) GetAccount(ctx context.Context, id string) (*Account, error) {
	acc := &Account{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, type, balance, user_id, created_at FROM accounts WHERE id = $1`,
		id).Scan(&acc.ID, &acc.Name, &acc.Type, &acc.Balance, &acc.UserID, &acc.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	return acc, nil
}

type TransactionRequest struct {
	ReferenceID string         `json:"reference_id"`
	Description string         `json:"description"`
	Entries     []EntryRequest `json:"entries"`
}

type EntryRequest struct {
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`    // Signed amount
	Direction string `json:"direction"` // Optional, helpful for validation
}

// RecordTransaction creates a transaction and its entries atomically.
// It also updates account balances.
func (r *Repository) RecordTransaction(ctx context.Context, req TransactionRequest) error {
	// 1. Validate Balance (Sum of amounts must be 0)
	var sum int64
	for _, e := range req.Entries {
		sum += e.Amount
	}
	if sum != 0 {
		return errors.New("transaction is not balanced (sum != 0)")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 2. Insert Transaction Record
	var transactionID string
	err = tx.QueryRowContext(ctx,
		`INSERT INTO transactions (reference_id, description) VALUES ($1, $2) RETURNING id`,
		req.ReferenceID, req.Description).Scan(&transactionID)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	// 3. Process Entries
	for _, e := range req.Entries {
		// Insert Entry
		_, err := tx.ExecContext(ctx,
			`INSERT INTO entries (transaction_id, account_id, amount, direction) VALUES ($1, $2, $3, $4)`,
			transactionID, e.AccountID, e.Amount, e.Direction)
		if err != nil {
			return fmt.Errorf("failed to create entry for account %s: %w", e.AccountID, err)
		}

		// Update Account Balance
		// Using simple addition update. Concurrency is handled by row-level lock implicitly acquired by UPDATE.
		// For high concurrency, might need 'FOR UPDATE' selects or specialized optimistic locking.
		// Here, generic atomic update is sufficient for MVP.
		res, err := tx.ExecContext(ctx,
			`UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
			e.Amount, e.AccountID)
		if err != nil {
			return fmt.Errorf("failed to update balance for account %s: %w", e.AccountID, err)
		}
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			return fmt.Errorf("account %s not found", e.AccountID)
		}
	}

	return tx.Commit()
}
