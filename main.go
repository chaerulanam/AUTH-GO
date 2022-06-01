package main

import (
	"auth/config"
	"auth/handlers"
	"auth/migration"
	"auth/repositories"
	"auth/services"
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := config.DB_USER + ":" + config.DB_PASS + "@tcp(" + config.SERVER_HOST + ":" + config.DB_PORT + ")/" + config.DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error")
	}

	migration.MigrateAll(db)

	userRepo := repositories.UserRepository(db)
	authService := services.AuthServ(userRepo)
	authHandler := handlers.AuthHandler(authService)

	e := echo.New()

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	e.Logger.Fatal(e.Start(":8888"))
}
