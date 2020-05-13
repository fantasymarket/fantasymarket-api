package api

import (
	v1 "fantasymarket/api/v1"

	"github.com/rs/zerolog/log"

	"fantasymarket/database"
	"fantasymarket/game"
	"fantasymarket/utils/config"

	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const addr = "localhost:5000"

// Start starts a new instance of our REST API
func Start(db *database.Service, game *game.Service, config *config.Config) {

	r := chi.NewRouter()

	// CORS Header
	cors := cors.New(cors.Options{
		AllowedOrigins:     []string{"https://fantasymarket.netlify.com/", "http://localhost:3000"},
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

	v1Handler := v1.NewAPIRouter(db, game, config)

	r.Mount("/v1", v1Handler)
	r.Mount("/", v1Handler)

	log.Info().Str("address", addr).Msg("successfully started the http server")
	http.ListenAndServe(addr, r)
}
