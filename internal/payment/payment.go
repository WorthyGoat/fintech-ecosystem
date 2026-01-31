package payment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// DB defines an interface for database operations.
type DB interface {
	QueryRowContext(ctx context.Context, query string, args ...any) Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// Row defines an interface for sql.Row.
type Row interface {
	Scan(dest ...any) error
}

// sqlDBWrapper wraps *sql.DB to satisfy the DB interface.
type sqlDBWrapper struct {
	*sql.DB
}

func (w *sqlDBWrapper) QueryRowContext(ctx context.Context, query string, args ...any) Row {
	return w.DB.QueryRowContext(ctx, query, args...)
}

// PaymentIntent represents a payment transaction intent.
type PaymentIntent struct {
	ID                   string    `json:"id"`
	Amount               int64     `json:"amount"` // In cents
	Currency             string    `json:"currency"`
	Status               string    `json:"status"` // requires_payment_method, succeeded, failed
	Description          string    `json:"description,omitempty"`
	UserID               string    `json:"user_id"`
	ApplicationFeeAmount int64     `json:"application_fee_amount,omitempty"`
	OnBehalfOf           string    `json:"on_behalf_of,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
}

// Repository handles database interactions for payments.
type Repository struct {
	db DB
}

// NewRepository creates a new instance of Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: &sqlDBWrapper{db}}
}

// NewTestRepository creates a repository with a custom DB interface for testing.
func NewTestRepository(db DB) *Repository {
	return &Repository{db: db}
}

// CreatePaymentIntent inserts a new payment intent.
func (r *Repository) CreatePaymentIntent(ctx context.Context, intent *PaymentIntent) error {
	if intent.Currency == "" {
		intent.Currency = "USD"
	}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO payment_intents (amount, currency, status, description, user_id, application_fee_amount, on_behalf_of) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`,
		intent.Amount, intent.Currency, intent.Status, intent.Description, intent.UserID, intent.ApplicationFeeAmount, intent.OnBehalfOf).
		Scan(&intent.ID, &intent.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment intent: %w", err)
	}
	return nil
}

// GetPaymentIntent retrieves a payment intent by ID.
func (r *Repository) GetPaymentIntent(ctx context.Context, id string) (*PaymentIntent, error) {
	var intent PaymentIntent
	err := r.db.QueryRowContext(ctx,
		"SELECT id, amount, currency, status, description, user_id, application_fee_amount, on_behalf_of, created_at FROM payment_intents WHERE id = $1",
		id).Scan(&intent.ID, &intent.Amount, &intent.Currency, &intent.Status, &intent.Description, &intent.UserID, &intent.ApplicationFeeAmount, &intent.OnBehalfOf, &intent.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}
	return &intent, nil
}

// UpdateStatus updates the status of a payment intent.
func (r *Repository) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE payment_intents SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}
	return nil
}

// IdempotencyRecord keys response.
type IdempotencyRecord struct {
	UserID       string
	Key          string
	ResponseBody string
	StatusCode   int
}

func (r *Repository) GetIdempotencyKey(ctx context.Context, userID, key string) (*IdempotencyRecord, error) {
	var rec IdempotencyRecord
	err := r.db.QueryRowContext(ctx, "SELECT user_id, key, response_body, status_code FROM idempotency_keys WHERE user_id = $1 AND key = $2", userID, key).
		Scan(&rec.UserID, &rec.Key, &rec.ResponseBody, &rec.StatusCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (r *Repository) SaveIdempotencyKey(ctx context.Context, userID, key string, statusCode int, body string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO idempotency_keys (user_id, key, response_body, status_code) VALUES ($1, $2, $3, $4)", userID, key, body, statusCode)
	return err
}
