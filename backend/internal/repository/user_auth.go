package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserAuthMethod struct {
	ID             uuid.UUID `db:"id"`
	UserId         uuid.UUID `db:"user_id"`
	Provider       string    `db:"provider"`
	ProviderUserID *string   `db:"provider_user_id"`
	PasswordHash   *string   `db:"password_hash"`
}

type UserAuthMethodRepositoryImp interface {
	GetByEmail(ctx context.Context, email string) (*UserAuthMethod, error)
	Create(ctx context.Context, m *UserAuthMethod) error
}

type UserAuthMethodRepository struct {
	db *pgxpool.Pool
}

func NewUserAuthMethodRepository(db *pgxpool.Pool) UserAuthMethodRepositoryImp {
	return &UserAuthMethodRepository{db: db}
}

func (r *UserAuthMethodRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*UserAuthMethod, error) {
	query := `SELECT uam.*
	FROM user_auth_methods uam
	JOIN users u ON u.user_id = uam.id
	WHERE u.email = $1 AND uam.provider = 'EMAIL_PASSWORD'
	`

	row := r.db.QueryRow(ctx, query, email)

	var m UserAuthMethod
	err := row.Scan(&m.ID,
		&m.UserId,
		&m.Provider,
		&m.ProviderUserID,
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
		id,
		user_id,
		provider,
		password_hash
	) VALUES (
		$1,
		$2,
		$3,
		$4
		
	)`

	_, err := r.db.Exec(ctx, query, m.ID, m.UserId, m.Provider, m.PasswordHash)
	if err != nil {
		return err
	}
	return nil
}
