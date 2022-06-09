package main

import (
	authHandlers "auth/V1/Auth/handlers"
	authRepo "auth/V1/Auth/repositories"
	authServices "auth/V1/Auth/services"
	"auth/config"
	"auth/helper"
	"auth/migration"
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

	userRepo := authRepo.UserRepository(db)
	authService := authServices.AuthServ(userRepo)
	authHandler := authHandlers.AuthHandler(authService)

	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/v1")

	v1.POST("/register", authHandler.Register)
	v1.POST("/login", authHandler.Login)

	// Configure middleware with the custom claims type

	p := v1.Group("/app")

	p.POST("/users", authHandler.GetUsers)
	p.POST("/group", authHandler.AddGroup, helper.IsAuth)
	p.POST("/permission", authHandler.AddPermission, helper.IsAuth)

	e.Logger.Fatal(e.Start(":8888"))
}
