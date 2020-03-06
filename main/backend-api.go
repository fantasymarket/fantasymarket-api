package main

import (
	"fantasymarket/requesthandler"
	"net/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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
	r.Get("/user/stats", requesthandler.UserStats)

	// GET news
	r.Get("/news", requesthandler.News) // Allow for query parameters

	// API-Routes
	r.Route("/stocks", func(r chi.Router) {

		// GET overview of all stocks
		r.Get("/", requesthandler.GetStockNumbers)

		// GET data from a specified stock
		r.Get("/{name}", requesthandler.GetStockDetails) // Allow for query parameters

		// POST data to make an order ( eg. SELL or BUY )
		r.Post("/orders", requesthandler.Orders)

	})

	http.ListenAndServe(":3000", r)
}
