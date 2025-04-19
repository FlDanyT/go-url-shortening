// unit тест

package jwt_test

import (
	"go/url-shortening/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {

	const email = "a@a.ru"

	jwtService := jwt.NewJWT("/2+XnmJGz1j3ehIVI/5P9kl+CghrE3DcS7rnT+qar5w=") // Передаем секретный токен
	token, err := jwtService.Create(jwt.JWTData{                             // Делаем токен
		Email: email,
	})

	if err != nil {
		t.Fatal(err)
	}

	isValid, data := jwtService.Parse(token) // Парсим токен получая gmail
	if !isValid {
		t.Fatal("Token is invalid")
	}

	if data.Email != email {
		t.Fatalf("Email %s not equal %s", data.Email, email)
	}

}
