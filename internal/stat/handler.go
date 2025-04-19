// Получение статистики по кликам

package stat

import (
	"go/url-shortening/configs"
	"time"

	"go/url-shortening/pkg/middleware"
	"go/url-shortening/pkg/res"
	"net/http"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandlerDeps struct {
	StatRepository *StatRepository
	Config         *configs.Config
}

type StatHandler struct { // Привязываем обработчики к структуре
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {

	handler := &StatHandler{

		StatRepository: deps.StatRepository,
	}

	router.Handle("GET /stat", middleware.IsAuthed(handler.GetStat(), deps.Config))

}

func (h *StatHandler) GetStat() http.HandlerFunc { // Группировки статистики по дням

	return func(w http.ResponseWriter, r *http.Request) {

		from, err := time.Parse("2006-01-02", r.URL.Query().Get("from")) // Парсим дату
		if err != nil {
			http.Error(w, "Invalid from param", http.StatusBadRequest)
			return
		}

		to, err := time.Parse("2006-01-02", r.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, "Invalid to param", http.StatusBadRequest)
			return
		}

		by := r.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid by param", http.StatusBadRequest)
			return
		}

		stats := h.StatRepository.GetStats(by, from, to)
		res.Json(w, stats, 200)

	}

}
