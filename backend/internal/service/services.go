package service

import (
	"github.com/C0deNe0/agromart/internal/lib/aws"
	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/repository"
)

type Services struct {
	User         *UserService
	Company      *CompanyService
	Product      *ProductService
	Auth         *AuthService
	RefreshToken *repository.RefreshTokenRepository
}

//later we can add the aws client directly here to the services which requires it

func NewServices(repo *repository.Repositories, tokenManager *utils.TokenManager, refreshTokenRepo *repository.RefreshTokenRepository, s3Client *aws.S3Service) *Services {

	CompanyService := NewCompanyService(repo.Company, repo.CompanyFollower)

	productService := NewProductService(repo.Product, repo.ProductImage, repo.ProductVariant, CompanyService.companyRepo, s3Client)

	return &Services{
		User:         NewUserService(repo.User),
		Company:      CompanyService,
		Product:      productService,
		Auth:         NewAuthService(repo.User, repo.UserAuthMethod, tokenManager, refreshTokenRepo),
		RefreshToken: refreshTokenRepo,
	}
}

// ❌ No secrets
// ❌ No config reads
// ❌ No infra creation
// ✅ Pure wiring
