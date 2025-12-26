package repository

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAuthMethod struct {
	model.Base
	UserId       uuid.UUID `db:"user_id"`
	AuthProvider string    `db:"auth_provider"`
	OAuthSub     *string   `db:"oauth_sub"`
	PasswordHash *string   `db:"password_hash"`
}

type UserAuthMethodRepository struct {
	db *pgxpool.Pool
}

func NewUserAuthMethodRepository(db *pgxpool.Pool) *UserAuthMethodRepository {
	return &UserAuthMethodRepository{db: db}
}

func (r *UserAuthMethodRepository) Create(
	ctx context.Context,
	m *UserAuthMethod,
) (*UserAuthMethod, error) {
	stmt := `
		INSERT INTO user_auth_methods (
			user_id,
			auth_provider,
			password_hash
		) VALUES (
			@user_id,
			@auth_provider,
			@password_hash
		
	)
		RETURNING *`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"user_id":       m.UserId,
		"auth_provider": m.AuthProvider,
		"password_hash": m.PasswordHash,
	})
	if err != nil {
		return nil, err
	}

	method, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserAuthMethod])
	if err != nil {
		return nil, err
	}
	return &method, nil
}

func (r *UserAuthMethodRepository) GetLocalByEmail(ctx context.Context, email string) (*UserAuthMethod, error) {
	query := ` SELECT
        uam.id,
        uam.user_id,
        uam.auth_provider,
        uam.oauth_sub,
        uam.password_hash
    FROM user_auth_methods uam
    JOIN users u ON u.id = uam.user_id
    WHERE u.email = $1
      AND uam.auth_provider = 'LOCAL'
	`

	var m UserAuthMethod

	err := r.db.QueryRow(ctx, query, email).Scan(
		&m.ID,
		&m.UserId,
		&m.AuthProvider,
		&m.OAuthSub,
		&m.PasswordHash,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *UserAuthMethodRepository) EnsureOAuth(ctx context.Context, userID uuid.UUID, provider string, sub string) (*UserAuthMethod, error) {
	query := `INSERT INTO user_auth_methods (
		user_id,
		auth_provider,
		oauth_sub
	) VALUES (
		@user_id,
		@auth_provider,
		@oauth_sub
	)
		ON CONFLICT (user_id, auth_provider)
		 DO UPDATE SET oauth_sub = EXCLUDED.oauth_sub
		RETURNING *
		`

	rows, err := r.db.Query(ctx, query, pgx.NamedArgs{
		"user_id":       userID,
		"auth_provider": provider,
		"oauth_sub":     sub,
	})
	if err != nil {
		return nil, err
	}
	method, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[UserAuthMethod])
	if err != nil {
		return nil, err
	}
	return &method, nil
}
