package repository

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	stmt := `INSERT INTO users (
			email,
			name,
			role,
			is_active,
			email_verified

		) VALUES (
			@email,
			@name,
			@role,
			@is_active,	
			@email_verified	
		)
		RETURNING *`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"email":          u.Email,
		"name":           u.Name,
		"role":           u.Role,
		"is_active":      u.IsActive,
		"email_verified": u.EmailVerified,
	})
	if err != nil {
		return nil, err
	}
	createdUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	var user user.User
	query := `SELECT *
			 FROM users 
			WHERE id = $1`

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
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]user.User, error) {

	stmt := `SELECT * FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) (*user.User, error) {
	stmt :=
		`
		UPDATE users
		SET
			name = @name,
			profile_image_url = @profile_image_url,
			phone = @phone,
			is_active = @is_active,
			email_verified = @email_verified,
			updated_at = NOW()
		WHERE id = @id
		RETURNING *
		
		`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"id":                u.ID,
		"name":              u.Name,
		"profile_image_url": u.ProfileImageURL,
		"phone":             u.Phone,
		"is_active":         u.IsActive,
		"email_verified":    u.EmailVerified,
	})
	if err != nil {
		return nil, err
	}
	updatedUser, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User

	query := `
			SELECT	*
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
