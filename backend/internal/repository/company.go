package repository

import "github.com/jackc/pgx/v5/pgxpool"

type CompanyRepositoryImp interface {
	Create()
	GetByID()
	List()
	Update()
}

type CompanyRepository struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) CompanyRepositoryImp {
	return &CompanyRepository{db: db}
}

func (r *CompanyRepository) Create() {

}

func (r *CompanyRepository) GetByID() {

}

func (r *CompanyRepository) List() {

}

func (r *CompanyRepository) Update() {

}
