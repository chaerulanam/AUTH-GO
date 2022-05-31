package main

import (
	"auth/config"
	"auth/migration"
	"fmt"

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

	fmt.Println("Hello")

}
