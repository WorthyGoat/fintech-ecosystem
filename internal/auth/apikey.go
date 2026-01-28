package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type APIKey struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	KeyPrefix    string     `json:"key_prefix"`
	KeyHash      string     `json:"-"` // Do not export
	TruncatedKey string     `json:"truncated_key"`
	Environment  string     `json:"environment"`
	Scopes       string     `json:"scopes"` // Space-separated scopes, "*" for all
	CreatedAt    time.Time  `json:"created_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
}

func (r *Repository) CreateAPIKey(ctx context.Context, key *APIKey) error {
	if key.Scopes == "" {
		key.Scopes = "*" // Default to all scopes
	}
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO api_keys (user_id, key_prefix, key_hash, truncated_key, environment, scopes)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		key.UserID, key.KeyPrefix, key.KeyHash, key.TruncatedKey, key.Environment, key.Scopes).
		Scan(&key.ID, &key.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create api key: %w", err)
	}
	return nil
}

func (r *Repository) GetAPIKeyByHash(ctx context.Context, hash string) (*APIKey, error) {
	var key APIKey
	var scopes sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, key_prefix, environment, scopes, revoked_at FROM api_keys WHERE key_hash = $1",
		hash).Scan(&key.ID, &key.UserID, &key.KeyPrefix, &key.Environment, &scopes, &key.RevokedAt)
	key.Scopes = scopes.String
	if key.Scopes == "" {
		key.Scopes = "*" // Default to all scopes for backwards compatibility
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get api key: %w", err)
	}
	return &key, nil
}
