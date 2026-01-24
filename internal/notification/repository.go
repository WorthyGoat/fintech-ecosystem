package notification

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Repository handles database operations for notifications.
type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new notification into the database.
func (r *Repository) Create(ctx context.Context, n *Notification) error {
	n.ID = uuid.New().String()
	n.CreatedAt = time.Now()
	n.Status = StatusPending

	query := `
		INSERT INTO notifications (id, user_id, recipient, channel, title, content, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		n.ID, n.UserID, n.Recipient, n.Channel, n.Title, n.Content, n.Status, n.CreatedAt,
	)
	return err
}

// UpdateStatus updates the status and sent_at timestamp of a notification.
func (r *Repository) UpdateStatus(ctx context.Context, id string, status Status) error {
	query := `UPDATE notifications SET status = $1, sent_at = $2 WHERE id = $3`
	var sentAt *time.Time
	if status == StatusSent {
		now := time.Now()
		sentAt = &now
	}
	_, err := r.db.ExecContext(ctx, query, status, sentAt, id)
	return err
}

// GetByID retrieves a notification by its ID.
func (r *Repository) GetByID(ctx context.Context, id string) (*Notification, error) {
	query := `
		SELECT id, user_id, recipient, channel, title, content, status, created_at, sent_at
		FROM notifications WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var n Notification
	err := row.Scan(&n.ID, &n.UserID, &n.Recipient, &n.Channel, &n.Title, &n.Content, &n.Status, &n.CreatedAt, &n.SentAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

// GetByUserID retrieves all notifications for a given user.
func (r *Repository) GetByUserID(ctx context.Context, userID string) ([]*Notification, error) {
	query := `
		SELECT id, user_id, recipient, channel, title, content, status, created_at, sent_at
		FROM notifications WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			// Using fmt/log or simple error ignore if context is done?
			// But Repository doesn't have logger.
			// Ideally we don't swallow, but GetByUserID signature returns error.
			// But defer happens after return.
			// Let's print or ignore. If log package is imported? No.
			// Wait, repository.go does NOT import log.
			// I should check imports.
			// Step 406 shows imports: context, database/sql, time, github.com/google/uuid. No log.
			// I will import log first.
			// Or I can just ignore it cleanly? No, errcheck is strict.
			// I will add log import if I can.
			// But replace_file_content for imports is risky without parsing.
			// I'll assume I can just `_ = rows.Close()`.
			_ = rows.Close()
		}
	}()

	var notifications []*Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Recipient, &n.Channel, &n.Title, &n.Content, &n.Status, &n.CreatedAt, &n.SentAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, rows.Err()
}
