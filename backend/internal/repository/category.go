package repository

import "github.com/jackc/pgx/v5/pgxpool"

type CategoryRepositoryImp interface {
	Create()
	GetByID()
	List()
	Update()
}

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepositoryImp {
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
