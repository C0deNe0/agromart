package service

import (
	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/repository"
)

type Services struct {
	User         *UserService
	Company      *CompanyService
	Product      *ProductService
	Auth         *AuthService
	RefreshToken *repository.RefreshTokenRepository
	Upload       *UploadService
}

//later we can add the aws client directly here to the services which requires it

func NewServices(repo *repository.Repositories, tokenManager *utils.TokenManager, refreshTokenRepo *repository.RefreshTokenRepository) *Services {

	// awsClient, err := aws.NewAWS(s)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create AWS client: %w", err)
	// }

	//---------------ERRORS NOT HANDLED PROPERLY HERE_____________________

	// tokenManager := utils.NewTokenManager("secret", "access")
	//____________________________________________________________________
	CompanyService := NewCompanyService(repo.Company, repo.CompanyFollower)
	productService := NewProductService(repo.Product, repo.Company, repo.Category)

	return &Services{
		User:         NewUserService(repo.User),
		Company:      CompanyService,
		Product:      productService,
		Auth:         NewAuthService(repo.User, repo.UserAuthMethod, tokenManager, refreshTokenRepo),
		RefreshToken: refreshTokenRepo,
		Upload:       NewUploadService(nil, productService),
	}
}

// ❌ No secrets
// ❌ No config reads
// ❌ No infra creation
// ✅ Pure wiring
