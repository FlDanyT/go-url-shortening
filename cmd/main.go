package main

import (
	"fmt"
	"go/url-shortening/pkg/event"
	"go/url-shortening/pkg/middleware"
	"net/http"

	"go/url-shortening/configs"
	"go/url-shortening/internal/auth"
	"go/url-shortening/internal/link"
	"go/url-shortening/internal/stat"
	"go/url-shortening/internal/user"
	"go/url-shortening/pkg/db"
)

func App() http.Handler {

	conf := configs.LoadConfig()
	db := db.NewDb(conf)

	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDeps{

		LinkRepository: linkRepository,
		Config:         conf,
		EventBus:       eventBus,
	})

	stat.NewStatHandler(router, stat.StatHandlerDeps{

		StatRepository: statRepository,
		Config:         conf,
	})

	go statService.AddClick() // Считаем клики по ссылкам

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)

}

func main() {

	app := App()
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()

}
