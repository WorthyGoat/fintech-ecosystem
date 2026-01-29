package auth

import (
	"context"
	"fmt"
	"time"
)

type Organization struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Membership struct {
	UserID    string    `json:"user_id"`
	OrgID     string    `json:"org_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Repository) CreateOrganization(ctx context.Context, name, domain string) (*Organization, error) {
	var org Organization
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO organizations (name, domain) VALUES ($1, $2) RETURNING id, name, domain, created_at",
		name, domain).Scan(&org.ID, &org.Name, &org.Domain, &org.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	return &org, nil
}

func (r *Repository) AddMember(ctx context.Context, userID, orgID, role string) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO memberships (user_id, org_id, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
		userID, orgID, role)
	return err
}

func (r *Repository) GetUserMemberships(ctx context.Context, userID string) ([]Membership, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT user_id, org_id, role, created_at FROM memberships WHERE user_id = $1",
		userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var memberships []Membership
	for rows.Next() {
		var m Membership
		if err := rows.Scan(&m.UserID, &m.OrgID, &m.Role, &m.CreatedAt); err != nil {
			return nil, err
		}
		memberships = append(memberships, m)
	}
	return memberships, nil
}
