package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTData struct {
	Email string
}

type JWT struct {
	Secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}

func (j *JWT) Create(data JWTData) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ // Делает новый jwt токен
		"email": data.Email,
	})

	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return s, nil

}

func (j *JWT) Parse(token string) (bool, *JWTData) { // Возвращаем email из jwt

	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { // Парсим jwt токен
		return []byte(j.Secret), nil
	})

	if err != nil {
		return false, nil
	}

	email := t.Claims.(jwt.MapClaims)["email"] // Декодирование json
	return t.Valid, &JWTData{
		Email: email.(string),
	}

}
