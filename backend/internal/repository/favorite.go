package repository

import "github.com/jackc/pgx/v5/pgxpool"

type FavoriteRepository struct {
	db *pgxpool.Pool
}

func NewFavoriteRepository(db *pgxpool.Pool) *FavoriteRepository {
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
