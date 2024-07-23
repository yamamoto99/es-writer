package handler

import (
	"es-app/controller"
	"es-app/middleware/auth"
	"es-app/middleware/cors"
	"es-app/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func NewRouter(
	ac controller.IAuthController,
	am auth.IAuthMiddleware,
	uc controller.IUserController,
	gc controller.IGenerateController,
) *echo.Echo {
	e := echo.New()
	e.Validator = validator.NewValidator()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.INFO)

	authGroup := e.Group("/auth")
	{
		authGroup.Use(cors.SetupAuthCORS())
		authGroup.POST("/signup", ac.SignUp)
		authGroup.POST("/checkEmail", ac.CheckEmail)
		authGroup.POST("/resendEmail", ac.ResendEmail)
		authGroup.POST("/login", ac.Login)
	}
	appGroup := e.Group("/app")
	{
		appGroup.Use(cors.SetupUserCORS())
		appGroup.Use(am.JwtMiddleware())
		userGroup := appGroup.Group("/profile")
		{
			userGroup.GET("/getProfile", uc.GetProfile)
			userGroup.PATCH("/updateProfile", uc.UpdateProfile)
		}
		generateGroup := appGroup.Group("/generate")
		{
			generateGroup.POST("/generateAnswers", gc.GenerateAnswers)
		}
	}
	return e
}
