package repository

import (
	"context"

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

func (r *CompanyRepository) Create(ctx context.Context, c *company.Company) (*company.Company, error) {
	stmt := `INSERT INTO companies (
		owner_id,
		 name,
		  description,
		   logo_url,
		    business_email,
		     business_phone,
		      city,
		       state,
		        pincode,
		         is_approved,
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
		FALSE,
		TRUE
		)
		RETURNING *`
	row, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"owner_id":       c.OwnerID,
		"name":           c.Name,
		"description":    c.Description,
		"logo_url":       c.LogoURL,
		"business_email": c.BusinessEmail,
		"business_phone": c.BusinessPhone,
		"city":           c.City,
		"state":          c.State,
		"pincode":        c.Pincode,
	})
	if err != nil {
		return nil, err
	}

	created, err := pgx.CollectOneRow(row, pgx.RowToStructByName[company.Company])
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	stmt := `SELECT * FROM companies WHERE id = $1`

	var c company.Company
	err := r.db.QueryRow(ctx, stmt, id).Scan(
		&c.ID,
		&c.OwnerID,
		&c.Name,
		&c.Description,
		&c.LogoURL,
		&c.BusinessEmail,
		&c.BusinessPhone,
		&c.City,
		&c.State,
		&c.Pincode,
		&c.IsApproved,
		&c.ApprovedBy,
		&c.ApprovedAt,
		&c.IsActive,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CompanyRepository) ListByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]company.Company, error) {
	stmt := `SELECT * 
	FROM companies 
	WHERE owner_id = $1
	ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, stmt, ownerID)
	if err != nil {
		return nil, err
	}

	companies, err := pgx.CollectRows(rows, pgx.RowToStructByName[company.Company])
	if err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *CompanyRepository) Update(ctx context.Context, c *company.Company) error {
	stmt := `UPDATE companies SET name = $2, description = $3, logo_url = $4, business_email = $5, business_phone = $6, city = $7, state = $8, pincode = $9 WHERE id = $1`
	_, err := r.db.Exec(ctx, stmt,
		c.ID,
		c.Name,
		c.Description,
		c.LogoURL,
		c.BusinessEmail,
		c.BusinessPhone,
		c.City,
		c.State,
		c.Pincode,
	)
	if err != nil {
		return err
	}
	return nil
}

func( r *CompanyRepository) SoftDelete(ctx context.Context, companyID uuid.UUID) error {
	stmt := `UPDATE companies
		SET is_active = FALSE
		updated_at = NOW()
		WHERE id = $1
		AND is_active = TRUE`
	_, err := r.db.Exec(ctx, stmt, companyID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CompanyRepository) ListPending(ctx context.Context) ([]company.Company, error) {
	stmt := `SELECT * FROM companies 
			WHERE is_approved = FALSE 
				AND is_active = TRUE
			ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}

	companies, err := pgx.CollectRows(rows, pgx.RowToStructByName[company.Company])
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func (r *CompanyRepository) Approve(ctx context.Context, adminID uuid.UUID, companyID uuid.UUID) error {
	query := `UPDATE companies
				SET
					is_approved = true,
					approved_by = $1,
					approved_at = NOW()
				WHERE id = $2`
	_, err := r.db.Exec(ctx, query, adminID, companyID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CompanyRepository) Reject(ctx context.Context, companyID uuid.UUID) error {
	query := `UPDATE companies
				SET is_active = false
				WHERE id = $1	
					AND is_approved = false`
	_, err := r.db.Exec(ctx, query, companyID)
	if err != nil {
		return err
	}
	return nil
}
