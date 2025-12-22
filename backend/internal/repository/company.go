package repository

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CompanyRepositoryImp interface {
	Create()
	GetByID()
	List()
	Update()
}

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) CompanyRepositoryImp {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create() {

}

func (r *CompanyRepository) GetByID() {

}

func (r *CompanyRepository) List() {

}

func (r *CompanyRepository) Update() {

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
