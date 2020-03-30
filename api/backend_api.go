package api

import (
	"fantasymarket/database"
	"fantasymarket/game"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var db *database.DatabaseService
var g *game.GameService

const addr = "localhost:42069"

// Start starts a new instance of our REST API
func Start(databaseService *database.DatabaseService, gameService *game.GameService) {
	db = databaseService
	g = gameService

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

	// Standalone GET Requests
	r.Get("/news", getEvents) // Allow for query parameters

	r.Get("/overview", getOverview) // Some stats for the dashboard
	// Top 2 Gainers / Top 2 Loosers
	// Maybe total + of all stock and things like that in the future

	r.Get("/time", getTime) // Current time on the server

	// API Routes
	r.Route("/stocks", func(r chi.Router) {

		r.Get("/", getAllStocks)

		r.Get("/{stockID}", getStockDetails)

		r.Post("/orders", addOrder)

	})

	r.Route("/orders", func(r chi.Router) {

		r.Get("/", orders)

		r.Get("/{orderID}", orders)

		r.Delete("/{orderID}", orders)
	})

	r.Route("/portfolio", func(r chi.Router) {

		r.Get("/", getPortfolio)

		r.Get("/{symbol}", getPortfolio)
	})

	fmt.Println("stated http server on " + addr + " :p")
	http.ListenAndServe(addr, r)
}
