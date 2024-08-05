package usecase

import (
	"es-app/infrastructure"
	"es-app/model"
	"es-app/repository"
	"os"

	"github.com/labstack/echo/v4"
)

type IAuthUsecase interface {
	SignUp(c echo.Context, signUpUser model.SignUpUser) (model.User, error)
	CheckEmail(c echo.Context, checkEmail model.CheckEmail) (bool, error)
	ResendEmail(c echo.Context, resendEmail model.ResendEmail) (bool, error)
	LogIn(c echo.Context, logInUser model.LoginUser) (model.LoginResponse, error)
	AccessToken(c echo.Context, accessToken string) (model.LoginUser, error)
	RefreshToken(c echo.Context, refreshToken string) (model.LoginResponse, model.LoginUser, error)
	LogOut(c echo.Context, accessToken string) error
	IsAlreadyRegisteredEmail(c echo.Context, email string) (bool, error)
}

type authUsecase struct {
	authRepo       repository.IAuthRepository
	infrastructure infrastructure.IIinfrastructure
	clientID       string
	jwtKeyURL      string
}

func NewAuthUsecase(authRepo repository.IAuthRepository, infrastructure infrastructure.IIinfrastructure) IAuthUsecase {
	return &authUsecase{
		authRepo:       authRepo,
		infrastructure: infrastructure,
		clientID:       os.Getenv("COGNITO_CLIENT_ID"),
		jwtKeyURL:      os.Getenv("TOKEN_KEY_URL"),
	}
}

func (u *authUsecase) IsAlreadyRegisteredEmail(c echo.Context, email string) (bool, error) {
	email, err := u.authRepo.FindByEmail(c, email)
	if err != nil {
		return false, err
	}
	return email != "", nil
}

func (au *authUsecase) SignUp(c echo.Context, signUpUser model.SignUpUser) (model.User, error) {
	user, err := au.infrastructure.SignUp(c, signUpUser)
	if err != nil {
		return model.User{}, err
	}

	if err := au.authRepo.CreateUser(c, user); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (au *authUsecase) CheckEmail(c echo.Context, checkEmail model.CheckEmail) (bool, error) {
	res, err := au.infrastructure.CheckEmail(c, checkEmail)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (au *authUsecase) ResendEmail(c echo.Context, resendEmail model.ResendEmail) (bool, error) {
	res, err := au.infrastructure.ResendEmail(c, resendEmail)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (au *authUsecase) LogIn(c echo.Context, logInUser model.LoginUser) (model.LoginResponse, error) {
	res, err := au.infrastructure.LogIn(c, logInUser)
	if err != nil {
		return model.LoginResponse{}, err
	}

	return res, nil
}

func (au *authUsecase) AccessToken(c echo.Context, accessToken string) (model.LoginUser, error) {
	res, err := au.infrastructure.ValidateToken(c, accessToken)
	if err != nil {
		return model.LoginUser{}, err
	}

	return res, nil
}

func (au *authUsecase) RefreshToken(c echo.Context, refreshToken string) (model.LoginResponse, model.LoginUser, error) {
	newToken, res, err := au.infrastructure.RefreshToken(c, refreshToken)
	if err != nil {
		return model.LoginResponse{}, model.LoginUser{}, err
	}

	return newToken, res, err
}

func (au *authUsecase) LogOut(c echo.Context, accessToken string) error {
	err := au.infrastructure.LogOut(c, accessToken)
	if err != nil {
		return err
	}
	return nil
}
