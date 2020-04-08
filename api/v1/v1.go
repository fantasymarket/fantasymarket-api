package v1

import (
	"fantasymarket/database"
	"fantasymarket/game"
	"net/http"

	"github.com/go-chi/chi"
)

// APIHandler holds the dependencies for http handlers
type APIHandler struct {
	DB   *database.Service
	Game *game.Service
}

// NewAPIRouter creates a new API HTTP handler
func NewAPIRouter(db *database.Service, game *game.Service) http.Handler {
	api := &APIHandler{
		DB:   db,
		Game: game,
	}

	r := chi.NewRouter()

	// Standalone GET Requests
	r.Get("/news", api.getEvents) // Allow for query parameters

	r.Get("/overview", api.getOverview) // Some stats for the dashboard
	// Top 2 Gainers / Top 2 Loosers
	// Maybe total + of all stock and things like that in the future

	r.Get("/time", api.getTime) // Current time on the server

	// API Routes
	r.Route("/stocks", func(r chi.Router) {

		r.Get("/", api.getAllStocks)

		r.Get("/{symbol}", api.getStockDetails)

		r.Post("/orders", api.addOrder)

	})

	r.Route("/orders", func(r chi.Router) {

		r.Get("/", api.orders)

		r.Post("/", api.addOrder)

		r.Get("/{orderID}", api.orders)

		r.Delete("/{orderID}", api.orders)
	})

	r.Route("/portfolio", func(r chi.Router) {

		r.Get("/", api.getPortfolio)

		r.Get("/{symbol}", api.getPortfolio)
	})

	return r
}
