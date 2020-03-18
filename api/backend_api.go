package api

import (
	"fantasymarket/database"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var db *database.DatabaseService

func Start(database *database.DatabaseService) {
	db = database
	r := chi.NewRouter()

	// CORS Header
	cors := cors.New(cors.Options{
		AllowedOrigins:     []string{"https://fantasymarket.netlify.com/"},
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
	r.Get("/news", GetEvents) // Allow for query parameters

	r.Get("/overview", GetOverview) // Some stats for the dashboard
	// Top 2 Gainers / Top 2 Loosers
	// Maybe total + of all stock and things like that in the future

	r.Get("/time", GetTime) // Current time on the server

	// API Routes
	r.Route("/stocks", func(r chi.Router) {

		r.Get("/", GetAllStocks)

		r.Get("/{stockID}", GetStockDetails)

		r.Post("/orders", AddOrder)

	})

	r.Route("/orders", func(r chi.Router) {

		r.Get("/", Orders)

		r.Get("/{orderID}", Orders)

		r.Delete("/{orderID}", Orders)
	})

	r.Route("/portfolio", func(r chi.Router) {

		r.Get("/", GetPortfolio)

		r.Get("/{symbol}", GetPortfolio)
	})

	http.ListenAndServe("localhost:42069", r)
}
