package controllerUtils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func SetSignupCookie(c echo.Context, name string, value string, expiration time.Duration) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = time.Now().Add(expiration)
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = http.SameSiteStrictMode

	c.SetCookie(cookie)
}

func SetLoginCookie(c echo.Context, idValue string, acValue string, refValue string) {
	c.SetCookie(&http.Cookie{
		Name:     "idToken",
		Value:    idValue,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "accessToken",
		Value:    acValue,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    refValue,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
}
