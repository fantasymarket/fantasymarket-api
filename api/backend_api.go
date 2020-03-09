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


	// GET user stats
	r.Get("/user/stats", UserStats)

	// GET news
	r.Get("/news", News) // Allow for query parameters

	// API-Routes
	r.Route("/stocks", func(r chi.Router) {

		// GET overview of all stocks
		r.Get("/", GetStockNumbers)

		// GET data from a specified stock
		r.Get("/{name}", GetStockDetails) // Allow for query parameters

		// POST data to make an order ( eg. SELL or BUY )
		r.Post("/orders", Orders)

	})

	http.ListenAndServe(":3000", r)
}