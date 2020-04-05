package api

import (
	v1 "fantasymarket/api/v1"
	"fantasymarket/database"
	"fantasymarket/game"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const addr = "localhost:42069"

// Start starts a new instance of our REST API
func Start(db *database.Service, game *game.Service) {

	r := chi.NewRouter()

	// CORS Header
	cors := cors.New(cors.Options{
		AllowedOrigins:     []string{"https://fantasymarket.netlify.com/", "http://" + addr},
		AllowedMethods:     []string{"GET", "POST", "PUT"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		MaxAge:             300,
		OptionsPassthrough: false,
		Debug:              false,
	})

	// Middleware
	r.Use(middleware.Logger, cors.Handler)

	v1Handler := v1.NewAPIRouter(db, game)

	r.Mount("/v1", v1Handler)
	r.Mount("/", v1Handler)

	fmt.Println("stated http server on " + addr + " :p")
	http.ListenAndServe(addr, r)
}
