// Unit test без базы данных

package auth_test

import (
	"go/url-shortening/internal/auth"
	"go/url-shortening/internal/user"
	"testing"
)

// Имитация базы данных
type MockUserRepository struct{}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "a@a.ru",
	}, nil
}

func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func TestRegisterSuccess(t *testing.T) {
	const initialEmail = "a@a.ru"
	authService := auth.NewAuthService(&MockUserRepository{}) // Подделка базы данных что бы не подключатся к ней

	email, err := authService.Register(initialEmail, "1", "Вася") // Делаем регистрацию
	if err != nil {
		t.Fatal(err)
	}

	if email != initialEmail {
		t.Fatalf("Email %s do not math %s", email, initialEmail)
	}

}
