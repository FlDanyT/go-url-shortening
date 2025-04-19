package res

import (
	"encoding/json"
	"net/http"
)

func Json(w http.ResponseWriter, data any, statusCode int) { // Ответ от api
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode) // Код ответа

	json.NewEncoder(w).Encode(data) // Возвращаем ответ от api
}
