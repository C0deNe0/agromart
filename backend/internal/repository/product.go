package repository

import (
	"context"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepositoryImp interface {
	Create(ctx context.Context, p *product.Product) (*product.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error)
	List(ctx context.Context, filter ProductFilter) (*model.PaginatedResponse[product.Product], error)
	Update(ctx context.Context, p *product.Product) (*product.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepositoryImp {
	return &ProductRepository{db: db}
}

type ProductFilter struct {
	CompanyID  *uuid.UUID
	CategoryID *uuid.UUID
	Search     *string
	Page       int
	Limit      int
	IsActive   *bool
}

func (r *ProductRepository) Create(ctx context.Context, p *product.Product) (*product.Product, error) {

	stmt := `
	INSERT INTO products (
			company_id,
			name,
			description,
			category_id,
			unit,
			origin,
			price_display,
			is_active
		)
		VALUES (
			@company_id,
			@name,
			@description,
			@category_id,
			@unit,
			@origin,
			@price_display,
			@is_active
		)
		RETURNING *`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"company_id":    p.CompanyID,
		"name":          p.Name,
		"description":   p.Description,
		"category_id":   p.CategoryID,
		"unit":          p.Unit,
		"origin":        p.Origin,
		"price_display": p.PriceDisplay,
		"is_active":     p.IsActive,
	})
	if err != nil {
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
		base += ` AND company_id=@company_id`
		args["company_id"] = *filter.CompanyID
	}

	if filter.CategoryID != nil {
		base += ` AND category_id=@category_id`
		args["category_id"] = *filter.CategoryID
	}

	if filter.Search != nil {
		base += ` AND name ILIKE @search`
		args["search"] = "%" + *filter.Search + "%"
	}

	if filter.IsActive != nil {
		base += ` AND is_active=@is_active`
		args["is_active"] = *filter.IsActive
	}

	//counting
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) `+base, args).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	//data
	stmt := `SELECT * ` + base + ` ORDER BY created_at DESC LIMIT @limit OFFSET @offset`

	args["limit"] = filter.Limit
	args["offset"] = (filter.Page - 1) * filter.Limit

	rows, err := r.db.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
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
			name=@name,
			description=@description,
			category_id=@category_id,
			unit=@unit,
			origin=@origin,
			price_display=@price_display,
			is_active=@is_active
		WHERE id=@id
		RETURNING *
	`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"id":            p.ID,
		"name":          p.Name,
		"description":   p.Description,
		"category_id":   p.CategoryID,
		"unit":          p.Unit,
		"origin":        p.Origin,
		"price_display": p.PriceDisplay,
		"is_active":     p.IsActive,
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
		UPDATE products SET is_active=false WHERE id=@id
	`
	_, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}
