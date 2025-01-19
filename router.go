package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"purchase-cart-service/dtos"
	"purchase-cart-service/errors"
	"purchase-cart-service/repositories"
	"purchase-cart-service/services"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

type Env struct {
	orders services.OrderService
}

func createRouterAndSetRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	env := &Env{
		orders: services.OrderService{
			OrderRepo:     repositories.OrderRepository{Db: db},
			OrderItemRepo: repositories.OrderItemRepository{Db: db},
			ProductRepo:   repositories.ProductRepository{Db: db},
		},
	}

	r.Post("/order", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		var body *dtos.OrderRequest
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			fmt.Println(err)
			return
		}

		order, err := env.orders.CreateOrder(ctx, body)
		if err != nil {
			if err == context.DeadlineExceeded {
				throwError(w, err.Error(), http.StatusRequestTimeout)
				return
			}

			handleError(w, err)
			return
		}

		writeJSON(w, http.StatusOK, order)
	})

	return r
}

func handleError(w http.ResponseWriter, err error) {
	if apiErr, ok := err.(*errors.APIError); ok {
		slog.Info("is an APIError")
		throwError(w, apiErr.Message, apiErr.StatusCode)
	} else {
		// For unexpected errors, return a generic response
		slog.Error(err.Error())
		throwError(w, "internal server error", http.StatusInternalServerError)
	}
}
