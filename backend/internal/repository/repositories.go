package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
	User             *UserRepository
	UserAuthMethod   *UserAuthMethodRepository
	Company          *CompanyRepository
	CompanyFollower  *CompanyFollowerRepository
	Category         *CategoryRepository
	Product          *ProductRepository
	RefreshToken     *RefreshTokenRepository
	Favorite         *FavoriteRepository
	SubscriptionPlan *SubscriptionPlanRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:             NewUserRepository(db),
		UserAuthMethod:   NewUserAuthMethodRepository(db),
		Company:          NewCompanyRepository(db),
		CompanyFollower:  NewCompanyFollowerRepository(db),
		Category:         NewCategoryRepository(db),
		Product:          NewProductRepository(db),
		RefreshToken:     NewRefreshTokenRepository(db),
		Favorite:         NewFavoriteRepository(db),
		SubscriptionPlan: NewSubscriptionPlanRepository(db),
	}
}
