package auth_test

import (
	"bytes"
	"encoding/json"
	"go/url-shortening/configs"
	"go/url-shortening/internal/auth"
	"go/url-shortening/internal/user"
	"go/url-shortening/pkg/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {

	database, mock, err := sqlmock.New() // Делаем mock базу
	if err != nil {

		return nil, nil, err
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{ // Подключимся к бд
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}

	userRepo := user.NewUserRepository(&db.Db{ // Передали созданную базу данных в репозиторий
		DB: gormDb,
	})

	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},

		AuthService: auth.NewAuthService(userRepo),
	}

	return &handler, mock, nil

}

func TestLoginHandlerSuccess(t *testing.T) {

	handler, mock, err := bootstrap()

	rows := sqlmock.NewRows([]string{"email", "password"}).AddRow("d@d.ru", "$2a$10$LdTOwaPiUiNYohAKTr6cGOppavk0X5Z/8CoxRrnargwzNEcTKo472")
	mock.ExpectQuery("SELECT").WillReturnRows(rows) // Если вызовется любой запрос должен вернуться rows
	if err != nil {

		t.Fatal(err)
		return

	}

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "d@d.ru",
		Password: "1",
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder() // В него записывать ответ

	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("got %d, expected %d", w.Code, 200)
	}

}

func TestRegisterHandlerSuccess(t *testing.T) {

	handler, mock, err := bootstrap()

	rows := sqlmock.NewRows([]string{"email", "password", "name"})

	mock.ExpectQuery("SELECT").WillReturnRows(rows) // Если вызовется любой запрос должен вернуться rows
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	if err != nil {

		t.Fatal(err)
		return

	}

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "d@d.ru",
		Password: "1",
		Name:     "Вася",
	})

	reader := bytes.NewReader(data)
	w := httptest.NewRecorder() // В него записывать ответ

	req := httptest.NewRequest(http.MethodPost, "/auth/register", reader)
	handler.Register()(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("got %d, expected %d", w.Code, 201)
	}

}
