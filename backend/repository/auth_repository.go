package repository

import (
	"es-app/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type IAuthRepository interface {
	CreateUser(c echo.Context, user model.User) error
	FindByEmail(c echo.Context, email string) (model.User, error)
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

func (r *authRepository) FindByEmail(c echo.Context, email string) (model.User, error) {
	var user model.User
	result := r.db.WithContext(c.Request().Context()).Select("email").Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.User{}, nil
		}
		return model.User{}, result.Error
	}
	return user, nil
}
