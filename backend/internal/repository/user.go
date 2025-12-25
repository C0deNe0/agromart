package repository

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryImp interface {
	Create(context.Context, *user.User) error
	GetByID(context.Context, uuid.UUID) (*user.User, error)
	List(context.Context) ([]user.User, error)
	Update(context.Context, *user.User) error
	GetByEmail(context.Context, string) (*user.User, error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryImp {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (
		email,
		name,
		role,
		is_active
	) VALUES (
		$1,
		$2,
		$3,
		$4
	)`

	_, err := r.db.Exec(ctx, query, u.Email, u.Name, u.Role, u.IsActive)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var user user.User
	query := `SELECT 
					id,
					email,
					name,
					profile_image_url,
					phone,
					role,
					is_active,
					email_verified,
					last_login_at,
					created_at,
					updated_at
	FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.ProfileImageURL,
		&user.Phone,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]user.User, error) {
	var users []user.User
	query := `SELECT * FROM users`

	err := r.db.QueryRow(ctx, query).Scan(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	query := `UPDATE users
	SET email = $2,
	name = $3,
	role = $4,
	is_active = $5
	WHERE id = $1`

	_, err := r.db.Exec(ctx, query, u.ID, u.Email, u.Name, u.Role, u.IsActive)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User

	query := `
		SELECT 
			id,
			email,
			name,
			profile_image_url,
			phone,
			role,
			is_active,
			email_verified,
			last_login_at,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.ProfileImageURL,
		&user.Phone,
		&user.Role,
		&user.IsActive,
		&user.EmailVerified,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
