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

	// API-Routes
	r.Route("/stocks", func(r chi.Router) {

		// Get - Requests
		r.Get("/", requesthandler.GetStockNumbers)
		r.Get("/portfolionumbers", requesthandler.GetPortfolioNumbers)

		// Post - Requests
		r.Post("/buy/{name}", requesthandler.BuyStock)
		r.Post("/sell/{name}", requesthandler.SellStock)

		// Put - Requests
		r.Put("/update/stocks", requesthandler.UpdateStocks)
	})

	http.ListenAndServe(":3000", r)
}
