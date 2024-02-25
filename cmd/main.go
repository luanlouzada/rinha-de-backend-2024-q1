package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/luanlouzada/rinha-de-backend-2024-q1/internal/handler"
)

func main() {
	r := chi.NewRouter()
	r.Post("/clientes/{ID}/transações", handler.HandleTransaction)
	r.Get("/clientes/[id]/extrato", HandleFunc)
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	http.ListenAndServe(":3000", r)
}
