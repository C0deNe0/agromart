package repository

import (
	"context"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyFollowerRepository struct {
	db *pgxpool.Pool
}

func NewCompanyFollowerRepository(db *pgxpool.Pool) *CompanyFollowerRepository {
	return &CompanyFollowerRepository{db: db}
}

func (r *CompanyFollowerRepository) Follow(ctx context.Context, companyID uuid.UUID, userID uuid.UUID) (*company.CompanyFollower, error) {
	stmt := `
	INSERT INTO company_followers (
		company_id,
		user_id
	)
	VALUES (
		@company_id,
		@user_id
	)
	ON CONFLICT (company_id, user_id) DO NOTHING
	RETURNING *`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"company_id": companyID,
		"user_id":    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to follow company:%w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[company.CompanyFollower])
	if err != nil {
		if err == pgx.ErrNoRows {
			return r.GetFollowStatus(ctx, companyID, userID)
		}
		return nil, fmt.Errorf("failed to collect row:%w", err)
	}

	return &row, nil
}

func (r *CompanyFollowerRepository) Unfollow(ctx context.Context, companyID uuid.UUID, userID uuid.UUID) error {
	stmt := `
	DELETE FROM company_followers
	WHERE company_id = @company_id AND user_id = @user_id`

	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"company_id": companyID,
		"user_id":    userID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow company:%w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user is not following company")
	}

	return nil
}

func (r *CompanyFollowerRepository) IsFollowing(ctx context.Context, companyID uuid.UUID, userID uuid.UUID) (bool, error) {
	stmt := `SELECT EXISTS (
		SELECT 1 
		FROM company_followers 
		WHERE company_id = @company_id AND user_id = @user_id
		)`
	var exists bool

	err := r.db.QueryRow(ctx, stmt, pgx.NamedArgs{
		"company_id": companyID,
		"user_id":    userID,
	}).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user is following company:%w", err)
	}

	return exists, nil
}

func (r *CompanyFollowerRepository) GetFollowStatus(ctx context.Context, companyID uuid.UUID, userID uuid.UUID) (*company.CompanyFollower, error) {
	stmt := `SELECT * FROM company_followers
		WHERE company_id = @company_id AND user_id = @user_id`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"company_id": companyID,
		"user_id":    userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get follow status:%w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[company.CompanyFollower])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to collect row:%w", err)
	}

	return &row, nil
}

func (r *CompanyFollowerRepository) GetFollowStatusBatch(ctx context.Context, companyIDs []uuid.UUID, userID uuid.UUID) (map[uuid.UUID]bool, error) {
	if len(companyIDs) == 0 {
		return make(map[uuid.UUID]bool), nil
	}

	stmt := `
        SELECT company_id
        FROM company_followers
        WHERE company_id = ANY(@company_ids) AND user_id = @user_id
    `

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"company_ids": companyIDs,
		"user_id":     userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get follow statuses: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]bool)
	for _, id := range companyIDs {
		result[id] = false
	}

	for rows.Next() {
		var companyID uuid.UUID
		if err := rows.Scan(&companyID); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result[companyID] = true
	}

	return result, nil
}

//DIRECT COPY PASTED

func (r *CompanyFollowerRepository) ListFollowers(ctx context.Context, companyID uuid.UUID, page, limit int) (*model.PaginatedResponse[company.CompanyFollowerResponse], error) {
	// Count total
	var total int
	countStmt := `SELECT COUNT(*) FROM company_followers WHERE company_id = @company_id`
	err := r.db.QueryRow(ctx, countStmt, pgx.NamedArgs{
		"company_id": companyID,
	}).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count followers: %w", err)
	}

	// Get data
	stmt := `
        SELECT 
            cf.id,
            cf.company_id,
            cf.user_id,
            u.name as user_name,
            u.email as user_email,
            cf.followed_at
        FROM company_followers cf
        JOIN users u ON cf.user_id = u.id
        WHERE cf.company_id = @company_id
        ORDER BY cf.followed_at DESC
        LIMIT @limit OFFSET @offset
    `

	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"company_id": companyID,
		"limit":      limit,
		"offset":     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list followers: %w", err)
	}
	defer rows.Close()

	var followers []company.CompanyFollowerResponse
	for rows.Next() {
		var f company.CompanyFollowerResponse
		err := rows.Scan(
			&f.ID,
			&f.CompanyID,
			&f.UserID,
			&f.UserName,
			&f.UserEmail,
			&f.FollowedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan follower: %w", err)
		}
		followers = append(followers, f)
	}

	return &model.PaginatedResponse[company.CompanyFollowerResponse]{
		Data:       followers,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + limit - 1) / limit,
	}, nil
}

func (r *CompanyFollowerRepository) ListFollowedCompanies(ctx context.Context, userID uuid.UUID, page, limit int) (*model.PaginatedResponse[company.Company], error) {
	// Count total
	var total int
	countStmt := `SELECT COUNT(*) FROM company_followers WHERE user_id = @user_id`
	err := r.db.QueryRow(ctx, countStmt, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count followed companies: %w", err)
	}

	// Get data
	stmt := `
        SELECT c.*
        FROM companies c
        JOIN company_followers cf ON c.id = cf.company_id
        WHERE cf.user_id = @user_id
        ORDER BY cf.followed_at DESC
        LIMIT @limit OFFSET @offset
    `

	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list followed companies: %w", err)
	}

	companies, err := pgx.CollectRows(rows, pgx.RowToStructByName[company.Company])
	if err != nil {
		return nil, fmt.Errorf("failed to collect companies: %w", err)
	}

	return &model.PaginatedResponse[company.Company]{
		Data:       companies,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (total + limit - 1) / limit,
	}, nil
}

func (r *CompanyFollowerRepository) CanViewProducts(ctx context.Context, companyID uuid.UUID, userID *uuid.UUID) (bool, error) {
	var visibility company.ProductVisibility
	var ownerID uuid.UUID
	var isApproved, isActive bool

	stmt := `
        SELECT product_visibility, owner_id, is_approved, is_active
        FROM companies 
        WHERE id = @company_id
    `

	err := r.db.QueryRow(ctx, stmt, pgx.NamedArgs{
		"company_id": companyID,
	}).Scan(&visibility, &ownerID, &isApproved, &isActive)

	if err != nil {
		return false, fmt.Errorf("failed to get company visibility: %w", err)
	}

	// Company must be approved and active
	if !isApproved || !isActive {
		// Only owner can view inactive/unapproved company products
		if userID != nil && ownerID == *userID {
			return true, nil
		}
		return false, nil
	}

	// PUBLIC - anyone can view
	if visibility == company.ProductVisibilityPublic {
		return true, nil
	}

	// Not authenticated - can't view non-public products
	if userID == nil {
		return false, nil
	}

	// Owner can always view
	if ownerID == *userID {
		return true, nil
	}

	// PRIVATE - only owner can view
	if visibility == company.ProductVisibilityPrivate {
		return false, nil
	}

	// FOLLOWERS_ONLY - check if user follows
	if visibility == company.ProductVisibilityFollowersOnly {
		return r.IsFollowing(ctx, companyID, *userID)
	}

	return false, nil
}
