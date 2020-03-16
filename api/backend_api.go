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
	r.Get("/news", GetNews) // Allow for query parameters

	r.Get("/overview", GetOverview) // Some stats for the dashboard

	r.Get("/time", GetTime) // Current time on the server

	// API Routes
	r.Route("/stocks", func(r chi.Router) {

		r.Get("/", GetStockNumbers)


		r.Get("/{symbol}", GetStockDetails)


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

	http.ListenAndServe(":3000", r)
}