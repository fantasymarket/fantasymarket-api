package api

import (
	v1 "fantasymarket/api/v1"

	"github.com/rs/zerolog/log"

	"fantasymarket/database"
	"fantasymarket/game"
	"fantasymarket/utils/config"

	"net/http"
	"github.com/go-chi/chi"
)

// Start starts a new instance of our REST API
func Start(db *database.Service, game *game.Service, config *config.Config) {

	r := chi.NewRouter()

	v1Handler := v1.NewAPIRouter(db, game, config)

	r.Mount("/v1", v1Handler)
	r.Mount("/", v1Handler)

	log.Info().Str("address", config.Port).Msg("successfully started the http server")
	http.ListenAndServe(":"+config.Port, r)
}
