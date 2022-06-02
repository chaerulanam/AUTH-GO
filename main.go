package main

import (
	"auth/config"
	"auth/handlers"
	"auth/helper"
	"auth/migration"
	"auth/repositories"
	"auth/services"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// func restricted(c echo.Context) error {
// 	user := c.Get("user").(*jwt.Token)
// 	claims := user.Claims.(*helper.JwtCustomClaims)
// 	username := claims.Username
// 	return c.String(http.StatusOK, "Welcome "+username+"!")
// }

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

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/register", authHandler.Register)
	e.POST("/login", authHandler.Login)

	// Configure middleware with the custom claims type

	e.POST("/user", authHandler.User, helper.IsAuth)

	e.Logger.Fatal(e.Start(":8888"))
}
