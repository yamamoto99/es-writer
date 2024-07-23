package controller

import (
	"es-app/model"
	"es-app/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	GetProfile(c echo.Context) error
	UpdateProfile(c echo.Context) error
}

type userController struct {
	userUsecase usecase.IUserUsecase
}

func NewUserController(userUsecase usecase.IUserUsecase) IUserController {
	return &userController{
		userUsecase: userUsecase,
	}
}

func (uc *userController) GetProfile(c echo.Context) error {
	userRes, err := uc.userUsecase.GetProfile(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, userRes)
}

func (uc *userController) UpdateProfile(c echo.Context) error {
	input := model.UserProfile{}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userRes, err := uc.userUsecase.UpdateProfile(c, input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, userRes)
}
