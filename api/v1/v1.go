package v1

import (
	"fantasymarket/database"
	"fantasymarket/game"
	"fantasymarket/utils/config"
	"fantasymarket/utils/http/middleware/jwt"
	"github.com/go-chi/chi"
	"net/http"
)

// APIHandler holds the dependencies for http handlers
type APIHandler struct {
	DB     *database.Service
	Game   *game.Service
	Config *config.Config
}

// NewAPIRouter creates a new API HTTP handler
func NewAPIRouter(db *database.Service, game *game.Service, config *config.Config) http.Handler {
	api := &APIHandler{
		DB:     db,
		Game:   game,
		Config: config,
	}

	r := chi.NewRouter()

	// Standalone GET Requests
	r.Get("/events", api.getEvents) // Allow for query parameters

	//r.Get("/overview", api.getOverview) // Some stats for the dashboard
	// Top 2 Gainers / Top 2 Losers
	// Maybe total + of all stock and things like that in the future

	r.Get("/time", api.getTime) // Current time on the server

	// API Routes
	r.Route("/stocks", func(r chi.Router) {

		r.Get("/", api.getAllStocks)

		r.Get("/{symbol}", api.getStockDetails)

	})

	r.Route("/orders", func(r chi.Router) {

		r.Get("/", api.ordersForUser)

		r.Post("/", api.addOrder)

		r.Get("/{orderID}", api.ordersID)

		r.Post("/fill/{orderID}", api.fillOrder)

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
