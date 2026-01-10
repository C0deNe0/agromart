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
	OwnerID        *uuid.UUID
	Search         *string
	ApprovalStatus *company.ApprovalStatus

	IsActive *bool
	Page     int
	Limit    int
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
		 gst_number,
		 pan_number,
		 product_visibility,
		 approval_status,
		 is_active
				  )
		VALUES (
		@owner_id,
		@name,
		@description,
		@logo_url,
		@business_email,
		@business_phone,
		@gst_number,
		@pan_number,
		@city,
		@state,
		@pincode,
		@product_visibility,
		@approval_status,
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
		"gst_number":         c.GSTNumber,
		"pan_number":         c.PANNumber,
		"product_visibility": c.ProductVisibility,
		"approval_status":    c.ApprovalStatus,
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

	if filter.ApprovalStatus != nil {
		base += ` AND is_approved = @is_approved`
		args["is_approved"] = *filter.ApprovalStatus
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
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
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
		gst_number = @gst_number,
		pan_number = @pan_number,
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
		"gst_number":         c.GSTNumber,
		"pan_number":         c.PANNumber,
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
		SET is_active = FALSE,
		updated_at= NOW() 
		WHERE id = @id`
	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("failed to delete company:%w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("company not found")
	}

	return nil
}

func (r *CompanyRepository) Approve(ctx context.Context, companyID uuid.UUID, adminID uuid.UUID, notes *string) error {
	query := `UPDATE companies
				SET
					approval_status = 'APPROVED',
            reviewed_by_id = @admin_id,
            reviewed_at = NOW(),
            rejection_reason = NULL,
            updated_at = NOW()
				WHERE id = @company_id AND approval_status = 'PENDING'`
	result, err := r.db.Exec(ctx, query, pgx.NamedArgs{
		"admin_id":   adminID,
		"company_id": companyID,
	})
	if err != nil {
		return fmt.Errorf("failed to approve company:%w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("company not found or not in pending status")
	}

	// loggin gthe approval histroy

	historyStmt := `
		INSERT INTO company_approval_history (company_id, action, performed_by_id, notes) 
		VALUES (@company_id, 'APPROVED', @admin_id,@notes)
	`
	_, err = r.db.Exec(ctx, historyStmt, pgx.StrictNamedArgs{
		"company_id": companyID,
		"admin_id":   adminID,
		"notes":      notes,
	})
	if err != nil {
		return fmt.Errorf("failed to log approval history: %w", err)
	}

	return nil
}

func (r *CompanyRepository) Reject(ctx context.Context, companyID, adminID uuid.UUID, reason string, notes *string) error {
	stmt := `
		UPDATE companies SET 
			approval_status = 'REJECTED',
			reviewed_by_id = @admin_id,
            reviewed_at = NOW(),
            rejection_reason = @reason,
            updated_at = NOW()
        WHERE id = @company_id AND approval_status = 'PENDING'
    `
	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"company_id": companyID,
		"admin_id":   adminID,
		"reason":     reason,
	})
	if err != nil {
		return fmt.Errorf("failed to reject comapny:%w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("company not found or not in pending status")
	}

	historyStmt := `
        INSERT INTO company_approval_history (company_id, action, performed_by_id, reason, notes)
        VALUES (@company_id, 'REJECTED', @admin_id, @reason, @notes)
    `
	_, err = r.db.Exec(ctx, historyStmt, pgx.NamedArgs{
		"company_id": companyID,
		"admin_id":   adminID,
		"reason":     reason,
		"notes":      notes,
	})
	if err != nil {
		return fmt.Errorf("failed to log rejection history: %w", err)
	}

	return nil

}

func (r *CompanyRepository) Resubmit(ctx context.Context, companyID, userID uuid.UUID) error {
	stmt := `
	 		UPDATE companies SET 
				approval_status = 'PENDING',
				 reviewed_by_id = NULL,
            reviewed_at = NULL,
            rejection_reason = NULL,
            submitted_at = NOW(),
            updated_at = NOW()
        WHERE id = @company_id AND approval_status = 'REJECTED'
	 	`

	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"company_id": companyID,
	})

	if err != nil {
		return fmt.Errorf("failed to resubmit company: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("comapny not found or not in rejected status")
	}

	historyStmt := `
        INSERT INTO company_approval_history (company_id, action, performed_by_id, notes)
        VALUES (@company_id, 'RESUBMITTED', @user_id, 'Company resubmitted after rejection')
    `
	_, err = r.db.Exec(ctx, historyStmt, pgx.NamedArgs{
		"company_id": companyID,
		"user_id":    userID,
	})
	if err != nil {
		return fmt.Errorf("failed to log resubmission history: %w", err)
	}

	return nil
}

func (r *CompanyRepository) GetApprovalHistory(ctx context.Context, companyID uuid.UUID) ([]company.CompanyApprovalHistory, error) {
	stmt := `
		SELECT * FROM company_approval_history
		WHERE company_id = @company_id
		ORDER BY created_at DESC
	
	`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"company_id": companyID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get approval history: %w", err)
	}

	history, err := pgx.CollectRows(rows, pgx.RowToStructByName[company.CompanyApprovalHistory])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return history, nil

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

func (r *CompanyRepository) GetApprovedCompanyByOwner(ctx context.Context, ownerID uuid.UUID) (*company.Company, error) {
	stmt := `
		SELECT * FROM companies
		WHERE owner_id = @owner_id
		AND approval_status = 'APPROVED'
		AND is_active = true
		LIMIT 1
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"owner_id": ownerID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get approved company: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[company.Company])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}

	return &row, nil
}

func (r *CompanyRepository) CountPendingApprovals(ctx context.Context) (int, error) {
	stmt := `SELECT COUNT(*) FROM companies WHERE approval_status = 'PENDING' AND is_active = true`

	var count int
	err := r.db.QueryRow(ctx, stmt).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending approvals: %w", err)
	}

	return count, nil
}
