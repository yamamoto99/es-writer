package controller

import (
	"encoding/json"
	"es-app/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IGenerateController interface {
	GenerateAnswers(c echo.Context) error
}

type generateController struct {
	generateUsecase usecase.IGenerateUsecase
}

func NewGenerateController(generateUsecase usecase.IGenerateUsecase) IGenerateController {
	return &generateController{
		generateUsecase: generateUsecase,
	}
}

func (gc *generateController) GenerateAnswers(c echo.Context) error {
	var html string
	err := json.NewDecoder(c.Request().Body).Decode(&html)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	answers, err := gc.generateUsecase.GenerateAnswers(c, html)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, answers)
}
