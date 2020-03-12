package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"net/http"
)


func main() {
	r := chi.NewRouter()

	// CORS Header
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"https://fantasymarket.netlify.com/"},
		//AllowOriginFunc:    nil,
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


	// GET news
	r.Get("/news", News) // Allow for query parameters

	r.Get("/overview", News) // Some stats for the dashboard

	r.Get("/time", News) // Current time on the server

	// POST requests
	r.Post("/orders", Orders) // Your currently open positions

	// API-Routes
	r.Route("/stocks", func(r chi.Router) {

		// GET overview of all stocks
		r.Get("/", GetStockNumbers)

		// GET data from a specified stock
		r.Get("/{symbol}", GetStockDetails) // Allow for query parameters

		// POST data to make an order ( eg. SELL or BUY )
		r.Post("/orders", Orders)

	})

	r.Route("/orders", func(r chi.Router) {

		r.Get("/", Orders)
		// GET specific order data
		r.Get("/{orderID}", Orders)

		r.Delete("/{orderID}", Orders)
	})

	r.Route("/portfolio", func(r chi.Router) {

		r.Get("/", Orders)

		r.Get("/{symbol}", Orders)
	})

	http.ListenAndServe(":3000", r)
}