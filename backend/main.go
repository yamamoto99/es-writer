package main

import (
	"es-app/controller"
	"es-app/db"
	"es-app/infrastructure"
	"es-app/middleware/auth"
	"es-app/repository"
	"es-app/handler"
	"es-app/usecase"
)

func main() {
	db := db.NewDB()
	authRepository := repository.NewAuthRepository(db)
	userRepository := repository.NewUserRepository(db)
	generateRepository := repository.NewGenerateRepository()
	infrastructure := infrastructure.NewInfrastructure()
	authUsecase := usecase.NewAuthUsecase(authRepository, infrastructure)
	userUsecase := usecase.NewUserUsecase(userRepository)
	generateUsecase := usecase.NewGenerateUsecase(generateRepository, userRepository)
	authController := controller.NewAuthController(authUsecase, userRepository)
	userController := controller.NewUserController(userUsecase)
	generateController := controller.NewGenerateController(generateUsecase)
	authMiddleware := auth.NewAuthMiddleware(infrastructure)
	e := handler.NewRouter(authController, authMiddleware, userController, generateController)
	e.Logger.Fatal(e.Start(":8080"))
}
