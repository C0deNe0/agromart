package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	User             UserRepositoryImp
	Company          CompanyRepositoryImp
	Product          ProductRepositoryImp
	Category         CategoryRepositoryImp
	Favorite         FavoriteRepositoryImp
	SubscriptionPlan SubscriptionPlanRepositoryImp
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		User:             NewUserRepository(db),
		Company:          NewCompanyRepository(db),
		Product:          NewProductRepository(db),
		Category:         NewCategoryRepository(db),
		Favorite:         NewFavoriteRepository(db),
		SubscriptionPlan: NewSubscriptionPlanRepository(db),
	}
}
