package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

// OAuthClient represents a registered OAuth client application.
type OAuthClient struct {
	ID               string    `json:"id"`
	ClientSecretHash string    `json:"-"`
	UserID           string    `json:"user_id"`
	Name             string    `json:"name"`
	IsPublic         bool      `json:"is_public"`
	CreatedAt        time.Time `json:"created_at"`
}

// OAuthToken represents an access token issued to a client.
type OAuthToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ClientID     string    `json:"client_id"`
	UserID       string    `json:"user_id,omitempty"`
	Scope        string    `json:"scope"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetClientByID retrieves a client by its ID.
func (r *Repository) GetClientByID(ctx context.Context, clientID string) (*OAuthClient, error) {
	var client OAuthClient
	err := r.db.QueryRowContext(ctx,
		"SELECT id, client_secret_hash, user_id, name, is_public, created_at FROM oauth_clients WHERE id = $1",
		clientID).Scan(&client.ID, &client.ClientSecretHash, &client.UserID, &client.Name, &client.IsPublic, &client.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}
	return &client, nil
}

// CreateOAuthToken stores a new access token.
func (r *Repository) CreateOAuthToken(ctx context.Context, token *OAuthToken) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO oauth_tokens (access_token, refresh_token, client_id, user_id, scope, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		token.AccessToken, token.RefreshToken, token.ClientID, token.UserID, token.Scope, token.ExpiresAt)

	if err != nil {
		return fmt.Errorf("failed to create token: %w", err)
	}
	return nil
}

// ValidateOAuthToken checks if an access token exists and is valid.
func (r *Repository) ValidateOAuthToken(ctx context.Context, accessToken string) (*OAuthToken, error) {
	var token OAuthToken
	err := r.db.QueryRowContext(ctx,
		"SELECT access_token, refresh_token, client_id, user_id, scope, expires_at, created_at FROM oauth_tokens WHERE access_token = $1",
		accessToken).Scan(&token.AccessToken, &token.RefreshToken, &token.ClientID, &token.UserID, &token.Scope, &token.ExpiresAt, &token.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	return &token, nil
}

// GenerateRandomString generates a secure random string (for tokens/secrets).
func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// HashString creates a SHA-256 hash of the input string.
func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

// AuthorizationCode represents an OAuth2 authorization code with PKCE.
type AuthorizationCode struct {
	Code                string    `json:"code"`
	ClientID            string    `json:"client_id"`
	UserID              string    `json:"user_id"`
	RedirectURI         string    `json:"redirect_uri"`
	Scope               string    `json:"scope"`
	CodeChallenge       string    `json:"code_challenge,omitempty"`
	CodeChallengeMethod string    `json:"code_challenge_method,omitempty"`
	ExpiresAt           time.Time `json:"expires_at"`
	Used                bool      `json:"used"`
	CreatedAt           time.Time `json:"created_at"`
}

// ClientRedirectURI represents an allowed redirect URI for an OAuth client.
type ClientRedirectURI struct {
	ID          string    `json:"id"`
	ClientID    string    `json:"client_id"`
	RedirectURI string    `json:"redirect_uri"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateAuthorizationCode stores a new authorization code.
func (r *Repository) CreateAuthorizationCode(ctx context.Context, code *AuthorizationCode) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO authorization_codes (code, client_id, user_id, redirect_uri, scope, code_challenge, code_challenge_method, expires_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		code.Code, code.ClientID, code.UserID, code.RedirectURI, code.Scope,
		code.CodeChallenge, code.CodeChallengeMethod, code.ExpiresAt)

	if err != nil {
		return fmt.Errorf("failed to create authorization code: %w", err)
	}
	return nil
}

// GetAuthorizationCode retrieves an authorization code by its value.
func (r *Repository) GetAuthorizationCode(ctx context.Context, code string) (*AuthorizationCode, error) {
	var authCode AuthorizationCode
	var codeChallenge, codeChallengeMethod sql.NullString

	err := r.db.QueryRowContext(ctx,
		`SELECT code, client_id, user_id, redirect_uri, scope, code_challenge, code_challenge_method, expires_at, used, created_at
		 FROM authorization_codes WHERE code = $1`,
		code).Scan(
		&authCode.Code, &authCode.ClientID, &authCode.UserID, &authCode.RedirectURI,
		&authCode.Scope, &codeChallenge, &codeChallengeMethod,
		&authCode.ExpiresAt, &authCode.Used, &authCode.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get authorization code: %w", err)
	}

	authCode.CodeChallenge = codeChallenge.String
	authCode.CodeChallengeMethod = codeChallengeMethod.String

	return &authCode, nil
}

// MarkAuthorizationCodeUsed marks an authorization code as used (one-time use).
func (r *Repository) MarkAuthorizationCodeUsed(ctx context.Context, code string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE authorization_codes SET used = TRUE WHERE code = $1 AND used = FALSE`,
		code)

	if err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("code already used or not found")
	}

	return nil
}

// CreateOAuthClient creates a new OAuth2 client.
func (r *Repository) CreateOAuthClient(ctx context.Context, client *OAuthClient) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO oauth_clients (id, client_secret_hash, user_id, name, is_public)
		 VALUES ($1, $2, $3, $4, $5)`,
		client.ID, client.ClientSecretHash, client.UserID, client.Name, client.IsPublic)

	if err != nil {
		return fmt.Errorf("failed to create oauth client: %w", err)
	}
	return nil
}

// AddRedirectURI adds a redirect URI to a client's allowlist.
func (r *Repository) AddRedirectURI(ctx context.Context, clientID, redirectURI string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO client_redirect_uris (client_id, redirect_uri) VALUES ($1, $2)
		 ON CONFLICT (client_id, redirect_uri) DO NOTHING`,
		clientID, redirectURI)

	if err != nil {
		return fmt.Errorf("failed to add redirect uri: %w", err)
	}
	return nil
}

// ValidateRedirectURI checks if redirect URI is allowed for the client.
func (r *Repository) ValidateRedirectURI(ctx context.Context, clientID, redirectURI string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM client_redirect_uris WHERE client_id = $1 AND redirect_uri = $2`,
		clientID, redirectURI).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to validate redirect uri: %w", err)
	}

	return count > 0, nil
}

// DeleteExpiredAuthorizationCodes removes expired authorization codes (cleanup).
func (r *Repository) DeleteExpiredAuthorizationCodes(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM authorization_codes WHERE expires_at < NOW()`)

	if err != nil {
		return fmt.Errorf("failed to delete expired codes: %w", err)
	}
	return nil
}
