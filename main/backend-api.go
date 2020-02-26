package main

import (
	// "fantasymarket/requesthandler"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func main() {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		//AllowOriginFunc:    nil,
		AllowedMethods:     []string{"GET"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		MaxAge:             300,
		OptionsPassthrough: false,
		Debug:              false,
	})

	r.Use(middleware.Logger, cors.Handler)

	// Serve the portfolio numbers
	// r.Get("/portfolionumbers", requesthandler.GetPortfolioNumbers())

	// Serve the stock numbers
	r.Route("/stocknumbers", func(r chi.Router) {
		// r.Get("/{name}", requesthandler.GetStockRequest)
	})

	http.ListenAndServe(":3000", r)
}
