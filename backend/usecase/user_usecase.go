package usecase

import (
	"es-app/model"
	"es-app/repository"

	"github.com/labstack/echo/v4"
)

type IUserUsecase interface {
	GetProfile(c echo.Context) (model.User, error)
	UpdateProfile(c echo.Context, input model.UserProfile) (model.User, error)
}

type userUsecase struct {
	userRepo repository.IUserRepository
}

func NewUserUsecase(userRepo repository.IUserRepository) IUserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (uu *userUsecase) GetProfile(c echo.Context) (model.User, error) {
	user, err := uu.userRepo.GetUser(c, c.Get("user_id").(string))
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (uu *userUsecase) UpdateProfile(c echo.Context, input model.UserProfile) (model.User, error) {
	userID := c.Get("user_id").(string)
	user, err := uu.userRepo.GetUser(c, userID)
	if err != nil {
		return model.User{}, err
	}

	user.Bio = input.Bio
	user.Experience = input.Experience
	user.Projects = input.Projects

	res, err := uu.userRepo.UpdateUser(c, userID, user)
	if err != nil {
		return model.User{}, err
	}

	return res, nil
}
