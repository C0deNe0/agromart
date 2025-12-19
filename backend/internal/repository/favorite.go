package repository

import "github.com/jackc/pgx/v5/pgxpool"

type FavoriteRepositoryImp interface {
	Create()
	GetByID()
	List()
	Update()
}

type FavoriteRepository struct {
	db *pgxpool.Pool
}

func NewFavoriteRepository(db *pgxpool.Pool) FavoriteRepositoryImp {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) Create() {

}

func (r *FavoriteRepository) GetByID() {

}

func (r *FavoriteRepository) List() {

}

func (r *FavoriteRepository) Update() {

}
