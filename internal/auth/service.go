package auth

import (
	"errors"
	"go/url-shortening/internal/user"
	"go/url-shortening/pkg/di"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository di.IUserRepository
}

func NewAuthService(UserRepository di.IUserRepository) *AuthService {

	return &AuthService{UserRepository: UserRepository}

}

func (service *AuthService) Login(email, password string) (string, error) {

	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser == nil { // Проверяем есть ли такой пользователь
		return "", errors.New(ErrorWrongCredetials)
	}

	err := bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(password)) // Расшифровываем пароль
	if err != nil {                                                                      // Проверяем есть ли такой пользователь
		return "", errors.New(ErrorWrongCredetials)
	}

	return existedUser.Email, nil

}

func (service *AuthService) Register(email, password, name string) (string, error) {

	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil { // Проверяем есть ли такой пользователь
		return "", errors.New(ErrUserExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Шифруем пароль
	if err != nil {
		return "", err
	}

	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	_, err = service.UserRepository.Create(user) // Делаем нового пользователя
	if err != nil {
		return "", err
	}

	return user.Email, nil

}
