package repository

import (
	"es-app/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateUser(c echo.Context, user model.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
	return &authRepository{db}
}

func (r *authRepository) CreateUser(c echo.Context, user model.User) error {
	return r.db.WithContext(c.Request().Context()).Create(&user).Error
}
