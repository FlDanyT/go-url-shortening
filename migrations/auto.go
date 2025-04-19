// go run migrations/auto.go

package main

import (
	"go/url-shortening/internal/link"
	"go/url-shortening/internal/stat"
	"go/url-shortening/internal/user"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() { // Подключаемся к базе данных и делаем миграции
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{}) // Делаем миграцию

}
