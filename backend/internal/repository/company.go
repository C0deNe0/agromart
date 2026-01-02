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

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{db: db}
}

type CompanyFilter struct {
	OwnerID    *uuid.UUID
	Search     *string
	IsApproved *bool
	IsActive   *bool
	Page       int
	Limit      int
}

func (r *CompanyRepository) Create(ctx context.Context, c *company.Company) (*company.Company, error) {
	stmt := `
	INSERT INTO companies (
		owner_id,
		 name,
		 description,
		 logo_url,
		 business_email,
		 business_phone,
		 city,
		 state,
		 pincode,
		 product_visibility,
		 is_active
				  )
		VALUES (
		@owner_id,
		@name,
		@description,
		@logo_url,
		@business_email,
		@business_phone,
		@city,
		@state,
		@pincode,
		@product_visibility,
		@is_active
		)
		RETURNING *`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"owner_id":           c.OwnerID,
		"name":               c.Name,
		"description":        c.Description,
		"logo_url":           c.LogoURL,
		"business_email":     c.BusinessEmail,
		"business_phone":     c.BusinessPhone,
		"city":               c.City,
		"state":              c.State,
		"pincode":            c.Pincode,
		"product_visibility": c.ProductVisibility,
		"is_active":          c.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create company:%w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[company.Company])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row:%w", err)
	}

	return &row, nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	stmt := `SELECT * FROM companies WHERE id = @id`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to get company by id:%w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[company.Company])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row:%w", err)
	}

	return &row, nil
}

func (r *CompanyRepository) List(ctx context.Context, filter CompanyFilter) (*model.PaginatedResponse[company.Company], error) {
	base := `FROM companies WHERE 1=1`
	args := pgx.NamedArgs{}

	if filter.OwnerID != nil {
		base += ` AND owner_id = @owner_id`
		args["owner_id"] = *filter.OwnerID
	}

	if filter.Search != nil {
		base += ` AND (name ILIKE @search OR description ILIKE @search)`
		args["search"] = "%" + *filter.Search + "%"
	}

	if filter.IsApproved != nil {
		base += ` AND is_approved = @is_approved`
		args["is_approved"] = *filter.IsApproved
	}

	if filter.IsActive != nil {
		base += ` AND is_active = @is_active`
		args["is_active"] = *filter.IsActive
	}

	//count
	var total int
	countStmt := `SELECT COUNT(*) ` + base
	if err := r.db.QueryRow(ctx, countStmt, args).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count companies:%w", err)
	}
	//get the data
	stmt := `SELECT * ` + base + ` ORDER BY created_at DESC LIMIT @limit OFFSET @offset`
	args["limit"] = filter.Limit
	args["offset"] = (filter.Page - 1) * filter.Limit

	rows, err := r.db.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to list companies:%w", err)
	}

	companies, err := pgx.CollectRows(rows, pgx.RowToStructByName[company.Company])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows:%w", err)
	}

	return &model.PaginatedResponse[company.Company]{
		Data:       companies,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      len(companies),
		TotalPages: (total + filter.Limit - 1) / filter.Limit,
	}, nil
}

func (r *CompanyRepository) Update(ctx context.Context, c *company.Company) (*company.Company, error) {
	stmt := `
	UPDATE companies
	SET
		name = @name,
		description = @description,
		logo_url = @logo_url,
		business_email = @business_email,
		business_phone = @business_phone,
		city = @city,
		state = @state,
		pincode = @pincode,
		product_visibility = @product_visibility,
		is_active = @is_active,
		updated_at = NOW()
	WHERE id = @id
	RETURNING *`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"id":                 c.ID,
		"name":               c.Name,
		"description":        c.Description,
		"logo_url":           c.LogoURL,
		"business_email":     c.BusinessEmail,
		"business_phone":     c.BusinessPhone,
		"city":               c.City,
		"state":              c.State,
		"pincode":            c.Pincode,
		"product_visibility": c.ProductVisibility,
		"is_active":          c.IsActive,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update company:%w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[company.Company])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row:%w", err)
	}

	return &row, nil
}

func (r *CompanyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	stmt := `UPDATE companies
		SET is_active = FALSE
		WHERE id = @id`
	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("failed to delete company:%w", err)
	}
	return nil
}

func (r *CompanyRepository) Approve(ctx context.Context, companyID uuid.UUID, adminID uuid.UUID) error {
	query := `UPDATE companies
				SET
					is_approved = true,
					approved_by = @admin_id,
					approved_at = NOW(),
					updated_at = NOW()
				WHERE id = @company_id`
	_, err := r.db.Exec(ctx, query, pgx.NamedArgs{
		"admin_id":   adminID,
		"company_id": companyID,
	})
	if err != nil {
		return fmt.Errorf("failed to approve company:%w", err)
	}
	return nil
}

func (r *CompanyRepository) GetByOwnerAndName(ctx context.Context, ownerID uuid.UUID, name string) (*company.Company, error) {
	stmt := `SELECT * FROM companies WHERE owner_id = @owner_id AND name = @name`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"owner_id": ownerID,
		"name":     name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get company by owner and name:%w", err)
	}
	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[company.Company])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to collect row:%w", err)
	}
	return &row, nil
}
