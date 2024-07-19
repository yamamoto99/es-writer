package model

import (
	"time"
)

type User struct {
	UserID     string    `json:"id" gorm:"gorm:unique not null"`
	Username   string    `json:"username" gorm:"unique not null"`
	Email      string    `json:"email"`
	Bio        string    `json:"bio"`
	Experience string    `json:"experience"`
	Projects   string    `json:"projects"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SignUpUser struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CheckEmail struct {
	Username         string `json:"username"`
	VerificationCode string `json:"verificationCode"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserProfile struct {
	Bio        string `json:"bio"`
	Experience string `json:"experience"`
	Projects   string `json:"projects"`
}
