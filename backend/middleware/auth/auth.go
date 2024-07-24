package auth

import (
	"es-app/infrastructure"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IAuthMiddleware interface {
	JwtMiddleware() echo.MiddlewareFunc
}

type authMiddleware struct {
	infrastructure infrastructure.IIinfrastructure
}

func NewAuthMiddleware(infrastructure infrastructure.IIinfrastructure) IAuthMiddleware {
	return &authMiddleware{infrastructure}
}

func (a *authMiddleware) JwtMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("accessToken")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
			}
			accessToken := cookie.Value
			if accessToken == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid token")
			}

			_, err = a.infrastructure.ValidateToken(c, accessToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
			}

			sub, err := a.infrastructure.GetUserID(c, accessToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token: "+err.Error())
			}

			c.Set("user_id", sub)

			return next(c)
		}
	}
}
