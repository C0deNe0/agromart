package repository

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, rt *auth.RefreshToken) (*auth.RefreshToken, error) {
	stmt := `
		INSERT INTO refresh_token (
			user_id,
			token_hash,
			user_agent,
			ip_address,
			expires_at
		) VALUES (
			@user_id,
			@token_hash,
			@user_agent,
			@ip_address,
			@expires_at
		)
		RETURNING *
	`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":    rt.UserID,
		"token_hash": rt.TokenHash,
		"user_agent": rt.UserAgent,
		"ip_address": rt.IPAddress,
		"expires_at": rt.ExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[auth.RefreshToken])
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *RefreshTokenRepository) FindValid(ctx context.Context, tokenHash string) (*auth.RefreshToken, error) {
	query := `SELECT *
	FROM refresh_token
	WHERE token_hash = $1
	AND revoked_at IS NULL
	AND expires_at > NOW()
	`
	var rt auth.RefreshToken
	err := r.db.QueryRow(ctx, query, tokenHash).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.UserAgent,
		&rt.IPAddress,
		&rt.ExpiresAt,
		&rt.RevokedAt,
		&rt.CreatedAt,
		&rt.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	stmt := `
		UPDATE refresh_token
		SET revoked_at = NOW()
		WHERE token_hash = @token_hash
		AND revoked_at IS NULL
	`
	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"token_hash": tokenHash,
	})
	return err
}

func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	stmt := `
		UPDATE refresh_token
		SET revoked_at = NOW()
		WHERE user_id = @user_id
		AND revoked_at IS NULL
	`
	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
	})
	return err
}

func (r *RefreshTokenRepository) IsValid(ctx context.Context, userID uuid.UUID, tokenHash string) (bool, error) {
	query := `SELECT *
	FROM refresh_token
	WHERE user_id = $1
	AND token_hash = $2
	AND revoked_at IS NULL
	AND expires_at > NOW()
	`
	var rt auth.RefreshToken
	err := r.db.QueryRow(ctx, query, userID, tokenHash).Scan(
		&rt.ID,
		&rt.UserID,
		&rt.TokenHash,
		&rt.UserAgent,
		&rt.IPAddress,
		&rt.ExpiresAt,
		&rt.RevokedAt,
		&rt.CreatedAt,
		&rt.UpdatedAt,
	)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ✔ multiple devices per user
// ✔ logout (single device)
// ✔ logout all devices
// ✔ refresh token rotation
// ✔ forced session revocation
// ✔ mobile-secure auth
