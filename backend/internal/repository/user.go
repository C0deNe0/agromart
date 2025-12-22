package repository

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryImp interface {
	Create(context.Context, *user.User) error
	GetByID(context.Context, uuid.UUID) (*user.UserResponse, error)
	List(context.Context) ([]user.UserResponse, error)
	Update(context.Context, *user.User) error
	GetByEmail(context.Context, string) (*user.UserResponse, error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryImp {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	query := `INSERT INTO users (
		id,
		email,
		name,
		role,
		is_active
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	)`

	_, err := r.db.Exec(ctx, query, u.ID, u.Email, u.Name, u.Role, u.IsActive)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.UserResponse, error) {
	var user user.UserResponse
	err := r.db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]user.UserResponse, error) {
	var users []user.UserResponse
	err := r.db.QueryRow(ctx, "SELECT * FROM users").Scan(&users)
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

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.UserResponse, error) {
	var user user.UserResponse
	err := r.db.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email).Scan(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
