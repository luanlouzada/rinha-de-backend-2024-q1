package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/luanlouzada/rinha-de-backend-2024-q1/internal"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "host=localhost:5432 user=admin password=123 dbname=rinha sslmode=disable")
	if err != nil {
		panic(err)
	}
	handler := &internal.Handler{
		DB: db,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/clientes/{id}/transações", handler.HandleTransaction)
	r.Get("/clientes/{id}/extrato", handler.HandleExtract)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	http.ListenAndServe(":8080", r)
}
