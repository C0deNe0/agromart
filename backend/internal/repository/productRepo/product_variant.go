package productRepo

import (
	"context"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductVariantRepository struct {
	db *pgxpool.Pool
}

func NewProductVariantRepository(db *pgxpool.Pool) *ProductVariantRepository {
	return &ProductVariantRepository{db: db}
}

func (r *ProductVariantRepository) Create(ctx context.Context, v *product.ProductVariant) (*product.ProductVariant, error) {
	stmt := `
        INSERT INTO product_variants (
            product_id,
            label,
            quantity_value,
            quantity_unit,
            price,
            stock_quantity,
            low_stock_threshold,
            is_available
        )
        VALUES (
            @product_id,
            @label,
            @quantity_value,
            @quantity_unit,
            @price,
            @stock_quantity,
            @low_stock_threshold,
            @is_available
        )
        RETURNING *
    `

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"product_id": v.ProductID,

		"label":               v.Label,
		"quantity_value":      v.QuantityValue,
		"quantity_unit":       v.QuantityUnit,
		"price":               v.Price,
		"stock_quantity":      v.StockQuantity,
		"low_stock_threshold": v.LowStockThreshold,
		"is_available":        v.IsAvailable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[product.ProductVariant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}

	return &row, nil
}

func (r *ProductVariantRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.ProductVariant, error) {
	stmt := `SELECT * FROM product_variants WHERE id = @id`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to get variant: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[product.ProductVariant])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("variant not found")
		}
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}

	return &row, nil
}

func (r *ProductVariantRepository) ListByProductID(ctx context.Context, productID uuid.UUID) ([]product.ProductVariant, error) {
	stmt := `
        SELECT * FROM product_variants
        WHERE product_id = @product_id
        ORDER BY price DESC
    `

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{"product_id": productID})
	if err != nil {
		return nil, fmt.Errorf("failed to list variants: %w", err)
	}

	variants, err := pgx.CollectRows(rows, pgx.RowToStructByName[product.ProductVariant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return variants, nil
}

func (r *ProductVariantRepository) ListByProductIDs(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID][]product.ProductVariant, error) {
	if len(productIDs) == 0 {
		return make(map[uuid.UUID][]product.ProductVariant), nil
	}

	stmt := `
        SELECT * FROM product_variants
        WHERE product_id = ANY(@product_ids)
        ORDER BY product_id, price ASC
    `

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{"product_ids": productIDs})
	if err != nil {
		return nil, fmt.Errorf("failed to list variants: %w", err)
	}

	variants, err := pgx.CollectRows(rows, pgx.RowToStructByName[product.ProductVariant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	// Group by product ID
	result := make(map[uuid.UUID][]product.ProductVariant)
	for _, v := range variants {
		result[v.ProductID] = append(result[v.ProductID], v)
	}

	return result, nil
}

func (r *ProductVariantRepository) Update(ctx context.Context, v *product.ProductVariant) (*product.ProductVariant, error) {
	stmt := `
        UPDATE product_variants SET
            label = @label,
            quantity_value = @quantity_value,
            quantity_unit = @quantity_unit,
            price = @price,
            stock_quantity = @stock_quantity,
            low_stock_threshold = @low_stock_threshold,
            is_available = @is_available,
            updated_at = NOW()
        WHERE id = @id
        RETURNING *
    `

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"id":                  v.ID,
		"label":               v.Label,
		"quantity_value":      v.QuantityValue,
		"quantity_unit":       v.QuantityUnit,
		"price":               v.Price,
		"stock_quantity":      v.StockQuantity,
		"low_stock_threshold": v.LowStockThreshold,
		"is_available":        v.IsAvailable,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[product.ProductVariant])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}

	return &row, nil
}

func (r *ProductVariantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	stmt := `DELETE FROM product_variants WHERE id = @id`

	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("variant not found")
	}

	return nil
}
