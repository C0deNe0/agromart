package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
	User             *UserRepository
	Company          *CompanyRepository
	CompanyFollower  *CompanyFollowerRepository
	Product          *ProductRepository
	Category         *CategoryRepository
	Favorite         *FavoriteRepository
	UserAuthMethod   *UserAuthMethodRepository
	SubscriptionPlan *SubscriptionPlanRepository
	RefreshToken     *RefreshTokenRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:             NewUserRepository(db),
		Company:          NewCompanyRepository(db),
		CompanyFollower:  NewCompanyFollowerRepository(db),
		Product:          NewProductRepository(db),
		Category:         NewCategoryRepository(db),
		Favorite:         NewFavoriteRepository(db),
		UserAuthMethod:   NewUserAuthMethodRepository(db),
		SubscriptionPlan: NewSubscriptionPlanRepository(db),
		RefreshToken:     NewRefreshTokenRepository(db),
	}
}
