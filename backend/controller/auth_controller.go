package controller

import (
	"es-app/controller/controllerUtils"
	"es-app/model"
	"es-app/usecase"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type IAuthController interface {
	SignUp(c echo.Context) error
	CheckEmail(c echo.Context) error
	ResendEmail(c echo.Context) error
	Login(c echo.Context) error
	LogOut(c echo.Context) error
}

type authController struct {
	authUsecase usecase.IAuthUsecase
}

func NewAuthController(authUsecase usecase.IAuthUsecase) IAuthController {
	return &authController{authUsecase}
}

func (ac *authController) SignUp(c echo.Context) error {
	signUpUser := model.SignUpUser{}
	if err := c.Bind(&signUpUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(signUpUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userRes, err := ac.authUsecase.SignUp(c, signUpUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	controllerUtils.SetSignupCookie(c, "username", signUpUser.Username, 10*time.Minute)

	return c.JSON(http.StatusCreated, userRes)
}

func (ac *authController) CheckEmail(c echo.Context) error {
	checkStruct := model.CheckEmail{}
	if err := c.Bind(&checkStruct); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(checkStruct); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	usernameCookie, err := c.Cookie("username")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	checkStruct.Username = usernameCookie.Value

	checkRes, err := ac.authUsecase.CheckEmail(c, checkStruct)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, checkRes)
}

func (ac *authController) ResendEmail(c echo.Context) error {
	resendEmail := model.ResendEmail{}

	usernameCookie, err := c.Cookie("username")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	resendEmail.Username = usernameCookie.Value

	resendRes, err := ac.authUsecase.ResendEmail(c, resendEmail)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resendRes)
}

func (ac *authController) Login(c echo.Context) error {
	loginUser := model.LoginUser{}
	accessToken, err := c.Cookie("accessToken")
	if err == nil {
		res, err := ac.authUsecase.AccessToken(c, accessToken.Value)
		if err == nil {
			c.Logger().Debug("🟡 Use access token")
			return c.JSON(http.StatusOK, res)
		}
	}

	refreshToken, err := c.Cookie("refreshToken")
	if err == nil {
		cookieRes, userRes, err := ac.authUsecase.RefreshToken(c, refreshToken.Value)
		if err == nil {
			c.Logger().Debug("🟡 Use refresh token")
			controllerUtils.SetLoginCookie(c, cookieRes.IDToken, cookieRes.AccessToken, cookieRes.RefreshToken)
			return c.JSON(http.StatusOK, userRes)
		}
	}

	if err := c.Bind(&loginUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(loginUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	loginRes, err := ac.authUsecase.LogIn(c, loginUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	c.Logger().Debug("🟡 Use username and password")

	controllerUtils.SetLoginCookie(c, loginRes.IDToken, loginRes.AccessToken, loginRes.RefreshToken)

	return c.JSON(http.StatusOK, loginUser)
}

func (ac *authController) LogOut(c echo.Context) error {
	accessToken, err := c.Cookie("accessToken")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Access token not found")
	}

	// Cognitoからのサインアウト
	err = ac.authUsecase.LogOut(c, accessToken.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// クライアント側でクッキーを削除するためのセット
	controllerUtils.ClearLoginCookie(c)

	c.Logger().Debug("🟡 Logged out")

	return c.NoContent(http.StatusOK)
}
