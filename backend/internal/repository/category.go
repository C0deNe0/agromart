package repository

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/category"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create() {

}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*category.Category, error) {
	return nil, nil
}

func (r *CategoryRepository) List() {

}

func (r *CategoryRepository) Update() {

}
