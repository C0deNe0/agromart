package productRepo

import (
	"context"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

type ProductFilter struct {
	CompanyID      *uuid.UUID
	CategoryID     *uuid.UUID
	Search         *string
	ApprovalStatus *company.ApprovalStatus
	IsActive       *bool
	Page           int
	Limit          int
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) (*product.Product, error) {

	stmt := `
	INSERT INTO products (
			company_id,
			category_id,
			name,
			description,
			unit,
			origin,
			base_price,
			approval_status,
			is_active
		)
		VALUES (
			@company_id,
			@category_id,
			@name,
			@description,
			@unit,
			@origin,
			@price,
			@approval_status,
			@is_active
		)
		RETURNING *`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"company_id":      p.CompanyID,
		"category_id":     p.CategoryID,
		"name":            p.Name,
		"description":     p.Description,
		"unit":            p.Unit,
		"origin":          p.Origin,
		"price":           p.BasePrice,
		"approval_status": p.ApprovalStatus,
		"is_active":       p.IsActive,
	})
	if err != nil {
		if pgErr, ok := err.(*pgx.ScanArgError); ok {
			return nil, fmt.Errorf("%s", pgErr.Err)
		}
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[product.Product])
	if err != nil {
		return nil, fmt.Errorf("failed to collect one row: %w", err)
	}

	return &row, nil

}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	stmt := `
		SELECT * FROM products
		WHERE id=@id
	`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[product.Product])
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	return &row, nil
}

func (r *ProductRepository) List(ctx context.Context, filter ProductFilter) (*model.PaginatedResponse[product.Product], error) {
	base := `FROM products WHERE 1=1`
	args := pgx.NamedArgs{}

	if filter.CompanyID != nil {
		base += ` AND company_id = @company_id`
		args["company_id"] = *filter.CompanyID
	}

	if filter.CategoryID != nil {
		base += ` AND category_id = @category_id`
		args["category_id"] = *filter.CategoryID
	}

	if filter.Search != nil {
		base += ` AND (name ILIKE @search OR description ILIKE @search)`
		args["search"] = "%" + *filter.Search + "%"
	}

	if filter.ApprovalStatus != nil {
		base += ` AND approval_status = @approval_status`
		args["approval_status"] = *filter.ApprovalStatus
	}

	if filter.IsActive != nil {
		base += ` AND is_active = @is_active`
		args["is_active"] = *filter.IsActive
	}

	// Count total
	var total int
	countStmt := `SELECT COUNT(*) ` + base
	if err := r.db.QueryRow(ctx, countStmt, args).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Get data
	stmt := `SELECT * ` + base + ` ORDER BY submitted_at DESC LIMIT @limit OFFSET @offset`
	args["limit"] = filter.Limit
	args["offset"] = (filter.Page - 1) * filter.Limit

	rows, err := r.db.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[product.Product])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return &model.PaginatedResponse[product.Product]{
		Data:       items,
		Total:      total,
		Page:       filter.Page,
		Limit:      filter.Limit,
		TotalPages: (total + filter.Limit - 1) / filter.Limit,
	}, nil
}
func (r *ProductRepository) Update(ctx context.Context, p *product.Product) (*product.Product, error) {

	stmt := `
		 UPDATE products SET 
		 	category_id=@category_id,
			name=@name,
			description=@description,
			unit=@unit,
			origin=@origin,
			base_price=@base_price,
			is_active=@is_active,
			updated_at= NOW()
		WHERE id=@id
		RETURNING *
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"id":          p.ID,
		"name":        p.Name,
		"description": p.Description,
		"category_id": p.CategoryID,
		"unit":        p.Unit,
		"origin":      p.Origin,
		"base_price":  p.BasePrice,
		"is_active":   p.IsActive,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[product.Product])
	if err != nil {
		return nil, fmt.Errorf("failed to collect one row: %w", err)
	}

	return &row, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	stmt := `
		UPDATE products SET is_active=false, updated_at=NOW() WHERE id=@id
	`
	res, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *ProductRepository) Approve(ctx context.Context, productID, adminID uuid.UUID, notes *string) error {
	stmt := `
		UPDATE products SET
			approval_status = 'APPROVED',
			reviewed_by_id = @admin_id,
			reviewed_at = NOW(),
			rejection_reason = NULL
		WHERE id = @product_id AND approval_status = 'PENDING'
	`
	res, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"product_id": productID,
		"admin_id":   adminID,
	})
	if err != nil {
		return fmt.Errorf("failed to approve product: %w", err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("product not found or not in pending status")
	}

	historyStmt := `
        INSERT INTO product_approval_history (product_id, action, performed_by_id, notes)
        VALUES (@product_id, 'APPROVED', @admin_id, @notes)
    `
	_, err = r.db.Exec(ctx, historyStmt, pgx.NamedArgs{
		"product_id": productID,
		"admin_id":   adminID,
		"notes":      notes,
	})
	if err != nil {
		return fmt.Errorf("failed to log approval history: %w", err)
	}

	return nil
}

func (r *ProductRepository) Reject(ctx context.Context, productID, adminID uuid.UUID, reason string, notes *string) error {
	stmt := `
        UPDATE products SET
            approval_status = 'REJECTED',
            reviewed_by_id = @admin_id,
            reviewed_at = NOW(),
            rejection_reason = @reason,
            updated_at = NOW()
        WHERE id = @product_id AND approval_status = 'PENDING'
    `

	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"product_id": productID,
		"admin_id":   adminID,
		"reason":     reason,
	})
	if err != nil {
		return fmt.Errorf("failed to reject product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found or not in pending status")
	}

	// Log rejection history
	historyStmt := `
        INSERT INTO product_approval_history (product_id, action, performed_by_id, reason, notes)
        VALUES (@product_id, 'REJECTED', @admin_id, @reason, @notes)
    `
	_, err = r.db.Exec(ctx, historyStmt, pgx.NamedArgs{
		"product_id": productID,
		"admin_id":   adminID,
		"reason":     reason,
		"notes":      notes,
	})
	if err != nil {
		return fmt.Errorf("failed to log rejection history: %w", err)
	}

	return nil
}

func (r *ProductRepository) Resubmit(ctx context.Context, productID, userID uuid.UUID) error {
	stmt := `
		UPDATE products SET
		 approval_status = 'PENDING',
            reviewed_by_id = NULL,
            reviewed_at = NULL,
            rejection_reason = NULL,
            submitted_at = NOW(),
            updated_at = NOW()
        WHERE id = @product_id AND approval_status = 'REJECTED'
    `

	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"product_id": productID,
	})
	if err != nil {
		return fmt.Errorf("failed to resubmit product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found or not in rejected status")
	}

	// Log resubmission
	historyStmt := `
        INSERT INTO product_approval_history (product_id, action, performed_by_id, notes)
        VALUES (@product_id, 'RESUBMITTED', @user_id, 'Product resubmitted after rejection')
    `
	_, err = r.db.Exec(ctx, historyStmt, pgx.NamedArgs{
		"product_id": productID,
		"user_id":    userID,
	})
	if err != nil {
		return fmt.Errorf("failed to log resubmission history: %w", err)
	}

	return nil
}

func (r *ProductRepository) GetApprovalHistory(ctx context.Context, productID uuid.UUID) ([]product.ProductApprovalHistory, error) {
	stmt := `
        SELECT * FROM product_approval_history
        WHERE product_id = @product_id
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{"product_id": productID})
	if err != nil {
		return nil, fmt.Errorf("failed to get approval history: %w", err)
	}

	history, err := pgx.CollectRows(rows, pgx.RowToStructByName[product.ProductApprovalHistory])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return history, nil
}

func (r *ProductRepository) CountPendingApprovals(ctx context.Context) (int, error) {
	stmt := `
		SELECT COUNT(*) FROM products WHERE approval_status = 'PENDING' 
		AND is_active = true
	`
	var count int
	err := r.db.QueryRow(ctx, stmt).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending approvals: %w", err)
	}
	return count, nil
}
