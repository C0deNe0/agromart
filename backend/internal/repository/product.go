package repository

import (
	"context"

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
		return nil, err
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[product.Product])
	if err != nil {
		return nil, err
	}

	return &row, nil

}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	return nil, nil
}

func (r *ProductRepository) List(ctx context.Context, filter ProductFilter) (*model.PaginatedResponse[product.Product], error) {
	return nil, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *product.Product) (*product.Product, error) {
	return nil, nil
}
