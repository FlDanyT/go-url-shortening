package middleware

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func Chain(middleware ...Middleware) Middleware { // Функция для добавления Middlewares

	return func(next http.Handler) http.Handler {

		for i := len(middleware) - 1; i >= 0; i-- {

			next = middleware[i](next)

		}

		return next

	}

}
