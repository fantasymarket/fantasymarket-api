package v1

import (
	"errors"
	"fantasymarket/database"
	"fantasymarket/game"
	"fantasymarket/utils/config"
	"fantasymarket/utils/http/middleware/jwt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
)

// APIHandler holds the dependencies for http handlers
type APIHandler struct {
	DB     *database.Service
	Game   *game.Service
	Config *config.Config
}

// Errors for the HTTP Handler
var (
	errFetchingData    = errors.New("error fetching data")
	errUpdatingOrder   = errors.New("error updating order")
	errDecoding        = errors.New("data could not be decoded")
	errDeletingOrder   = errors.New("order could not be deleted")
	errUserNotFound    = errors.New("could not find user")
	errParsingPassword = errors.New("could not parse password")
	errParsingUsername = errors.New("could not parse username")
	errGeneratingToken = errors.New("could not generate token")
	errCreatingAccount = errors.New("error creating new user account")
	errInvalidLogin    = errors.New("could not login user")
	errStockNotFound   = errors.New("could not find stock")
)

// NewAPIRouter creates a new API HTTP handler
func NewAPIRouter(db *database.Service, game *game.Service, config *config.Config) http.Handler {
	api := &APIHandler{
		DB:     db,
		Game:   game,
		Config: config,
	}

	r := chi.NewRouter()

	// CORS Header
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://fantasymarket.netlify.app", "https://develop--fantasymarket.netlify.app", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Middleware
	r.Use(corsConfig.Handler)
	r.Use(middleware.Logger)

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
			r.Post("/login", api.loginUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(jwt.Middleware(api.Config.TokenSecret, false))
			r.Get("/", api.getSelf)
			r.Post("/", api.updateSelf)
		})
	})

	return r
}
