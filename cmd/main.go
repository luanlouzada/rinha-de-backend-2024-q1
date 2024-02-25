package main

import (
	"net/http"

	"github.com/luanlouzada/rinha-de-backend-2024-q1/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Post("/clientes/{ID}/transações", handler.HandleTransaction)
	r.Get("/clientes/[id]/extrato")
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	http.ListenAndServe(":3000", r)
}
