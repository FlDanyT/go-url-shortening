// Оброботка запросов
package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler { // Логгер

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		wrapper := &WrapperWriter{ // Делаем обертку с StatusCode
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)
		log.Println(wrapper.StatusCode, r.Method, r.URL.Path, time.Since(start)) // Выводим логи

	})

}
