package repository

import (
	"github.com/C0deNe0/agromart/internal/repository/productRepo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	User             *UserRepository
	UserAuthMethod   *UserAuthMethodRepository
	Company          *CompanyRepository
	CompanyFollower  *CompanyFollowerRepository
	Category         *CategoryRepository
	Product          *productRepo.ProductRepository
	ProductImage     *productRepo.ProductImageRepository
	ProductVariant   *productRepo.ProductVariantRepository
	RefreshToken     *RefreshTokenRepository
	// Favorite         *FavoriteRepository
	SubscriptionPlan *SubscriptionPlanRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:             NewUserRepository(db),
		UserAuthMethod:   NewUserAuthMethodRepository(db),
		Company:          NewCompanyRepository(db),
		CompanyFollower:  NewCompanyFollowerRepository(db),
		Category:         NewCategoryRepository(db),
		Product:          productRepo.NewProductRepository(db),
		ProductImage:     productRepo.NewProductImageRepository(db),
		ProductVariant:   productRepo.NewProductVariantRepository(db),
		RefreshToken:     NewRefreshTokenRepository(db),
		// Favorite:         NewFavoriteRepository(db),
		SubscriptionPlan: NewSubscriptionPlanRepository(db),
	}
}
