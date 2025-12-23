package repository

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyRepositoryImp interface {
	Create(context.Context, *company.Company) (*company.Company, error)
	GetByID(context.Context, uuid.UUID) (*company.Company, error)
	ListByOwnerID(context.Context, uuid.UUID) ([]company.Company, error)
	Update(context.Context, *company.Company) (*company.Company, error)
	ListPending(context.Context) ([]company.Company, error)
	Approve(context.Context, uuid.UUID, uuid.UUID) error
	Reject(context.Context, uuid.UUID, uuid.UUID) error
}

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) CompanyRepositoryImp {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create(ctx context.Context, c *company.Company) (*company.Company, error) {
	stmt := `INSERT INTO companies (owner_id, name, description, logo_url, business_email, business_phone, city, state, pincode, is_approved, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, false, true)
			RETURNING *`
	row := r.db.QueryRow(ctx, stmt,
		c.OwnerID,
		c.Name,
		c.Description,
		c.LogoURL,
		c.BusinessEmail,
		c.BusinessPhone,
		c.City,
		c.State,
		c.Pincode,
	)

	err := row.Scan(
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

	return c, nil
}

func (r *CompanyRepository) GetByID(ctx context.Context, id uuid.UUID) (*company.Company, error) {
	stmt := `SELECT * FROM companies WHERE id = $1`
	row := r.db.QueryRow(ctx, stmt, id)
	var c company.Company
	err := row.Scan(
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
	stmt := `SELECT * FROM companies WHERE owner_id = $1`
	rows, err := r.db.Query(ctx, stmt, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []company.Company
	for rows.Next() {
		var c company.Company
		err := rows.Scan(
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
		companies = append(companies, c)
	}
	return companies, nil
}

func (r *CompanyRepository) Update(ctx context.Context, c *company.Company) (*company.Company, error) {
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
		return nil, err
	}
	return c, nil
}

func (r *CompanyRepository) ListPending(ctx context.Context) ([]company.Company, error) {
	query := `SELECT * FROM companies WHERE is_approved = false 
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companies []company.Company
	for rows.Next() {
		var company company.Company
		if err := rows.Scan(
			&company.ID,
			&company.OwnerID,
			&company.Name,
			&company.Description,
			&company.LogoURL,
			&company.BusinessEmail,
			&company.BusinessPhone,
			&company.City,
			&company.State,
			&company.Pincode,
			&company.IsApproved,
			&company.ApprovedBy,
			&company.ApprovedAt,
			&company.IsActive,
			&company.CreatedAt,
			&company.UpdatedAt,
		); err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}
	return companies, nil
}

func (r *CompanyRepository) Approve(ctx context.Context, adminID uuid.UUID, companyID uuid.UUID) error {
	query := `UPDATE companies
		 SET is_approved = true,
		  approved_by = $1,
		   approved_at = NOW()
		    WHERE id = $2`
	_, err := r.db.Exec(ctx, query, adminID, companyID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CompanyRepository) Reject(ctx context.Context, adminID uuid.UUID, companyID uuid.UUID) error {
	query := `UPDATE companies
			 SET is_approved = false,
			  approved_by = $1,
			   approved_at = NOW()
			    WHERE id = $2`
	_, err := r.db.Exec(ctx, query, adminID, companyID)
	if err != nil {
		return err
	}
	return nil
}
