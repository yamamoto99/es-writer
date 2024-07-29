package repository

import (
	"es-app/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUser(c echo.Context, id string) (model.User, error)
	UpdateUser(c echo.Context, id string, input model.User) (model.User, error)
	FindByEmail(c echo.Context, email string) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUser(c echo.Context, id string) (model.User, error) {
	var user model.User
	result := r.db.WithContext(c.Request().Context()).First(&user, "user_id = ?", id)
	if result.Error != nil {
		return model.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return model.User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) UpdateUser(c echo.Context, id string, input model.User) (model.User, error) {
	var user model.User
	result := r.db.Model(&user).Where("user_id = ?", id).Updates(input).WithContext(c.Request().Context())
	if result.Error != nil {
		return model.User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) FindByEmail(c echo.Context, email string) (model.User, error) {
	var user model.User
	result := r.db.WithContext(c.Request().Context()).First(&user, "email = ?", email)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return model.User{}, nil
		}
		return model.User{}, result.Error
	}
	return user, nil
}
