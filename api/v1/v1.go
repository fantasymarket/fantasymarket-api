package v1

import (
	"fantasymarket/database"
	"fantasymarket/game"
	"fantasymarket/utils/config"
	"fantasymarket/utils/http/middleware/jwt"
	"net/http"
	"errors"
	"github.com/go-chi/chi"
)

// APIHandler holds the dependencies for http handlers
type APIHandler struct {
	DB     *database.Service
	Game   *game.Service
	Config *config.Config
}

// Errors for the HTTP Handler
var (
	fetchingError = errors.New("error fetching data")
	orderUpdateError = errors.New("error updating error")
	decodingError = errors.New("data could not be decoded")
	orderDeletionError = errors.New("order could not be deleted")
	userNotFoundError = errors.New("could not find user")
	passwordError = errors.New("could not parse password")
	usernameError = errors.New("could not parse username")
	tokenError = errors.New("could not generate token")
	accountError = errors.New("error creating new user account")
	loginError = errors.New("could not login user")
)

// NewAPIRouter creates a new API HTTP handler
func NewAPIRouter(db *database.Service, game *game.Service, config *config.Config) http.Handler {
	api := &APIHandler{
		DB:     db,
		Game:   game,
		Config: config,
	}

	r := chi.NewRouter()

	r.Get("/events", api.getEvents)


	r.Get("/time", api.getTime)


	r.Route("/stocks", func(r chi.Router) {

		r.Get("/", api.getAllStocks)

		r.Get("/{symbol}", api.getStockDetails)

	})

	r.Route("/orders", func(r chi.Router) {

		r.Get("/", api.ordersForUser)

		r.Post("/", api.addOrder)

		r.Get("/{orderID}", api.getOrdersID)

		r.Delete("/{orderID}", api.deleteOrder)
	})

	r.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwt.Middleware(api.Config.TokenSecret, true))
			r.Get("/{username}", api.getUser)
			r.Put("/", api.createUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(jwt.Middleware(api.Config.TokenSecret, false))
			r.Get("/", api.getSelf)
			r.Post("/", api.updateSelf)
		})
	})

	return r
}
