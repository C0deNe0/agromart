package repository

import "github.com/jackc/pgx/v5/pgxpool"


type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create() {

}

func (r *CategoryRepository) GetByID() {

}

func (r *CategoryRepository) List() {

}

func (r *CategoryRepository) Update() {

}
