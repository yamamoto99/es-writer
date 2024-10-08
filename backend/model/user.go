package model

import (
	"time"
)

type User struct {
	UserID         string    `json:"id" gorm:"gorm:unique not null"`
	Username       string    `json:"username" gorm:"unique not null"`
	Email          string    `json:"email"`
	WorkExperience string    `json:"workExperience"`
	Skills         string    `json:"skills"`
	SelfPR         string    `json:"selfPR"`
	FutureGoals    string    `json:"futureGoals"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SignUpUser struct {
	Username  string    `json:"username" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CheckEmail struct {
	Username         string `json:"username"`
	VerificationCode string `json:"verificationCode" validate:"required"`
}

type ResendEmail struct {
	Username string `json:"username"`
}

type LoginUser struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserProfile struct {
	WorkExperience string `json:"workExperience"`
	Skills         string `json:"skills"`
	SelfPR         string `json:"selfPR"`
	FutureGoals    string `json:"futureGoals"`
}
