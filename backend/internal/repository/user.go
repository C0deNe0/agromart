package repository

import "github.com/jackc/pgx/v5/pgxpool"

type UserRepositoryImp interface {
	Create()
	GetByID()
	List()
	Update()
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepositoryImp {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create() {

}

func (r *UserRepository) GetByID() {

}

func (r *UserRepository) List() {

}

func (r *UserRepository) Update() {

}
