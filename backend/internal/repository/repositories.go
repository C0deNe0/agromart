package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
	User             UserRepositoryImp
	Company          CompanyRepositoryImp
	Product          ProductRepositoryImp
	Category         CategoryRepositoryImp
	Favorite         FavoriteRepositoryImp
	UserAuthMethod   UserAuthMethodRepositoryImp
	SubscriptionPlan SubscriptionPlanRepositoryImp
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:             NewUserRepository(db),
		Company:          NewCompanyRepository(db),
		Product:          NewProductRepository(db),
		Category:         NewCategoryRepository(db),
		Favorite:         NewFavoriteRepository(db),
		UserAuthMethod:   NewUserAuthMethodRepository(db),
		SubscriptionPlan: NewSubscriptionPlanRepository(db),
	}
}
