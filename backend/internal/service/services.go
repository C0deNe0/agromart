package service

import (
	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/repository"
)

type Services struct {
	User    *UserService
	Company *CompanyService
	Product *ProductService
	Auth    *AuthService
}

//later we can add the aws client directly here to the services which requires it

func NewServices(repo *repository.Repositories) *Services {

	// awsClient, err := aws.NewAWS(s)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create AWS client: %w", err)
	// }



	//---------------ERRORS NOT HANDLED PROPERLY HERE_____________________

	tokenManager := utils.NewTokenManager("secret", "access")
	//____________________________________________________________________

	
	return &Services{
		User:    NewUserService(repo.User),
		Company: NewCompanyService(repo.Company),
		Product: NewProductService(repo.Product, repo.Company),
		Auth:    NewAuthService(repo.User, repo.UserAuthMethod, tokenManager),
	}
}
