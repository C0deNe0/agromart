package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAuthMethod struct {
	ID           uuid.UUID `db:"id"`
	UserId       uuid.UUID `db:"user_id"`
	AuthProvider string    `db:"auth_provider"`
	OAuthSub     *string   `db:"oauth_sub"`
	PasswordHash *string   `db:"password_hash"`
}

type UserAuthMethodRepositoryImp interface {
	GetByEmail(ctx context.Context, email string) (*UserAuthMethod, error)
	Create(ctx context.Context, m *UserAuthMethod) error
	EnsureOAuth(ctx context.Context, userID uuid.UUID, provider string, providerUserID string) error
}

type UserAuthMethodRepository struct {
	db *pgxpool.Pool
}

func NewUserAuthMethodRepository(db *pgxpool.Pool) UserAuthMethodRepositoryImp {
	return &UserAuthMethodRepository{db: db}
}

func (r *UserAuthMethodRepository) GetByEmail(ctx context.Context, email string) (*UserAuthMethod, error) {
	query := `SELECT uam.*
	FROM user_auth_methods uam
	JOIN users u ON u.id = uam.user_id
	WHERE u.email = $1 AND uam.auth_provider = 'LOCAL'
	`

	row := r.db.QueryRow(ctx, query, email)

	var m UserAuthMethod
	err := row.Scan(&m.ID,
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

func (r *UserAuthMethodRepository) Create(
	ctx context.Context,
	m *UserAuthMethod,
) error {
	query := `INSERT INTO user_auth_methods (
		user_id,
		auth_provider,
		password_hash
	) VALUES (
		$1,
		'LOCAL',
		$2
		
	)`

	_, err := r.db.Exec(ctx, query, m.UserId, m.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserAuthMethodRepository) EnsureOAuth(ctx context.Context, userID uuid.UUID, provider string, providerUserID string) error {
	query := `INSERT INTO user_auth_methods (
		user_id,
		auth_provider,
		oauth_sub
	) VALUES (
		$1,
		'GOOGLE',
		$2,
		
		
	)`

	_, err := r.db.Exec(ctx, query, userID, providerUserID)
	if err != nil {
		return err
	}
	return nil
}
