package main

import (
	"bytes"
	"encoding/json"
	"go/url-shortening/internal/auth"
	"go/url-shortening/internal/user"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB { // Доступ к базе данных для тестов

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db

}

func initData(db *gorm.DB) { // Создаем пользователя для тестовой таблички

	db.Create(&user.User{
		Email:    "a2@a.ru",
		Password: "$2a$10$LdTOwaPiUiNYohAKTr6cGOppavk0X5Z/8CoxRrnargwzNEcTKo472",
		Name:     "Daniil",
	})

}

func removeData(db *gorm.DB) { // Удаляет созданного пользователя

	db.Unscoped().
		Where("email = ?", "a2@a.ru").
		Delete(&user.User{})

}

func TestLoginSuccess(t *testing.T) { // Проверяет работает ли авторизация правильно

	// Prepare
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{ // Передаем данные
		Email:    "a2@a.ru",
		Password: "1",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data)) // Делаем запрос и смотрим будут ли ошибки
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body) // Получаем токен
	if err != nil {
		t.Fatal(err)
	}

	var resData auth.LoginResponse

	err = json.Unmarshal(body, &resData) // Раскодируем токен
	if err != nil {
		t.Fatal(err)
	}

	if resData.Token == "" {
		t.Fatal("Token empty")
	}

	removeData(db)

}

func TestLoginFail(t *testing.T) { // Проверяет работает ли авторизация правильно

	// Prepare
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{ // Передаем реальные данные
		Email:    "a2@a.ru",
		Password: "1",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data)) // Делаем запрос и смотрим будут ли ошибки
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	removeData(db)

}
