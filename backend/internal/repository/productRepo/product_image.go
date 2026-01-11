package productRepo

import (
	"context"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductImageRepository struct {
	db *pgxpool.Pool
}

func NewProductImageRepository(db *pgxpool.Pool) *ProductImageRepository {
	return &ProductImageRepository{db: db}
}

func (r *ProductImageRepository) Create(ctx context.Context, img *product.ProductImage) (*product.ProductImage, error) {
	stmt := `
		INSERT INTO products_images (
			product_id,
			image_url,
			s3_key,
			is_primary
		)
		VALUES (@product_id, @image_url, @s3_key, @is_primary)
        RETURNING *
    `

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"product_id": img.ProductID,
		"image_url":  img.ImageURL,
		"s3_key":     img.S3Key,
		"is_primary": img.IsPrimary,
		// "display_order": img.DisplayOrder,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create product image: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[product.ProductImage])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}

	return &row, nil
}

func (r *ProductImageRepository) GetByID(ctx context.Context, id uuid.UUID) (*product.ProductImage, error) {
	stmt := `SELECT * FROM product_images WHERE id = @id`

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to get image: %w", err)
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[product.ProductImage])
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("image not found")
		}
		return nil, fmt.Errorf("failed to collect row: %w", err)
	}

	return &row, nil
}
func (r *ProductImageRepository) ListByProductID(ctx context.Context, productID uuid.UUID) ([]product.ProductImage, error) {
	stmt := `
        SELECT * FROM product_images
        WHERE product_id = @product_id
        ORDER BY is_primary DESC
    `

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"product_id": productID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	images, err := pgx.CollectRows(rows, pgx.RowToStructByName[product.ProductImage])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	return images, nil
}

func (r *ProductImageRepository) ListByProductIDs(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID][]product.ProductImage, error) {
	if len(productIDs) == 0 {
		return make(map[uuid.UUID][]product.ProductImage), nil
	}

	stmt := `
        SELECT * FROM product_images
        WHERE product_id = ANY(@product_ids)
        ORDER BY product_id, is_primary DESC
    `

	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{"product_ids": productIDs})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	images, err := pgx.CollectRows(rows, pgx.RowToStructByName[product.ProductImage])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}

	// Group by product ID
	result := make(map[uuid.UUID][]product.ProductImage)
	for _, img := range images {
		result[img.ProductID] = append(result[img.ProductID], img)
	}

	return result, nil
}

func (r *ProductImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	stmt := `DELETE FROM product_images WHERE id = @id`

	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}

func (r *ProductImageRepository) SetPrimary(ctx context.Context, productID, imageID uuid.UUID) error {
	// Start transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Unset all primary images for this product
	_, err = tx.Exec(ctx, `
        UPDATE product_images 
        SET is_primary = false 
        WHERE product_id = $1
    `, productID)
	if err != nil {
		return fmt.Errorf("failed to unset primary images: %w", err)
	}

	// Set new primary image
	result, err := tx.Exec(ctx, `
        UPDATE product_images 
        SET is_primary = true 
        WHERE id = $1 AND product_id = $2
    `, imageID, productID)
	if err != nil {
		return fmt.Errorf("failed to set primary image: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("image not found")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

