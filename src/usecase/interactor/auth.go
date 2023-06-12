package interactor

import (
	"context"
	"strings"
	"todo-service/graph/model"
	"todo-service/src/infrastructure/authentication"
	"todo-service/src/models"
	"todo-service/src/usecase/repository"
	"todo-service/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authInteractor struct {
	AuthRepository  repository.AuthRepository
	UserRepository  repository.UserRepository
	jwtConfigurator authentication.JwtConfigurator
}

type AuthInteractor interface {
	SignUp(ctx context.Context, input model.NewUser) (*model.SignUpResult, error)
	SignIn(ctx context.Context, email string, password string) (*model.SignInResult, error)
	ValidateJwtToken(bearerToken string) error
}

func NewAuthInteractor(
	r repository.AuthRepository, p repository.UserRepository, jc authentication.JwtConfigurator) AuthInteractor {
	return &authInteractor{r, p, jc}
}

func (ai *authInteractor) SignUp(ctx context.Context, input model.NewUser) (*model.SignUpResult, error) {
	// Check Email
	userDb, err := ai.UserRepository.GetByEmail(input.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, models.ErrInternalServerError
	}

	if userDb != nil {
		return nil, models.ErrUserEmailAlreadyExists
	}

	input.Password = utils.HashPwd(input.Password)

	user := models.User{
		ID:       uuid.New(),
		Name:     input.Name,
		Email:    strings.ToLower(input.Email),
		Password: input.Password,
	}

	_, err = ai.UserRepository.Create(user)
	if err != nil {
		return nil, err
	}

	return &model.SignUpResult{
		IsCreated: true,
	}, nil

}

func (ai *authInteractor) SignIn(ctx context.Context, email string, password string) (*model.SignInResult, error) {
	getUser, err := ai.UserRepository.GetByEmail(email)
	if err != nil {
		// if user not found
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrUserEmailNotFound
		}
		return nil, err
	}

	if err := utils.ComparePwd(getUser.Password, password); err != nil {
		return nil, models.ErrUserPasswordIsInvalid
	}

	token, err := ai.jwtConfigurator.CreateTokenPair(ctx, getUser.ID)
	if err != nil {
		return nil, err
	}

	return &model.SignInResult{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}

func (ai *authInteractor) ValidateJwtToken(bearerToken string) error {

	_, err := ai.jwtConfigurator.ValidateJwtToken(bearerToken)
	if err != nil {
		return models.ErrInvalidAccessToken
	}

	return nil
}
