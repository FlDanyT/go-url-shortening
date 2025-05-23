package link

import (
	"fmt"
	"go/url-shortening/configs"

	"go/url-shortening/pkg/event"
	"go/url-shortening/pkg/middleware"
	"go/url-shortening/pkg/req"
	"go/url-shortening/pkg/res"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository

	Config *configs.Config

	EventBus *event.EventBus
}

type LinkHandler struct { // Привязываем обработчики к структуре
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {

	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}

	router.Handle("POST /link", middleware.IsAuthed(handler.Create(), deps.Config))
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.Delete(), deps.Config))
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.Handle("GET /link", middleware.IsAuthed(handler.GetAll(), deps.Config))

}

func (handler *LinkHandler) Create() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		body, err := req.HandleBody[LinkCreateRequest](&w, r)
		if err != nil {
			return
		}

		link := NewLink(body.Url) // Делаем ссылку

		for {

			existedLink, _ := handler.LinkRepository.GetByHash(link.Hash) // Проверяем уникальность хеша

			if existedLink == nil {
				break
			}
			link.GenerateHash()

		}

		createdLink, err := handler.LinkRepository.Create(link) // Отправляем на запись в бд
		if err != nil {

			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		res.Json(w, createdLink, 201)

	}

}

func (handler *LinkHandler) Update() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		email, ok := r.Context().Value(middleware.ContextEmailKey).(string) // Достаем значение из ContextEmailKey
		if ok {
			fmt.Println(email)
		}

		body, err := req.HandleBody[LinkUpdateRequest](&w, r) // Получем переданое значение
		if err != nil {
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link, err := handler.LinkRepository.Update(&Link{ // Делаем новое значение с новыми данными
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, link, 201)

	}

}

func (handler *LinkHandler) Delete() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.GetById(uint(id))
		if err != nil { // Если нету записи в бд возвращаем ошибку
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = handler.LinkRepository.Delete(uint(id))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, nil, 200)

	}

}

func (handler *LinkHandler) GoTo() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		hash := r.PathValue("hash")                         // Получим значение из запросы
		link, err := handler.LinkRepository.GetByHash(hash) // Ищет хеш в бд

		if err != nil {

			http.Error(w, err.Error(), http.StatusNotFound)
			return

		}

		go handler.EventBus.Publush(event.Event{
			Type: event.EventLinkVisited,
			Data: link.ID,
		})
		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect) // Отправляем страницу сайта из ссылки

	}

}

func (handler *LinkHandler) GetAll() http.HandlerFunc { // Получаем все ссылки

	return func(w http.ResponseWriter, r *http.Request) {

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusBadRequest)
			return
		}

		// Получаем данные
		links := handler.LinkRepository.GetAll(limit, offset)
		count := handler.LinkRepository.Count()

		res.Json(w, GetAllLinksResponse{ // Выдаем ответ
			Links: links,
			Count: count,
		}, 200)

	}

}
