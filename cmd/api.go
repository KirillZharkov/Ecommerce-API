package main

import (
	"log"
	"net/http"
	"time"

	repo "github.com/KirillZharkov/Ecommerce-API/internal/adapters/postgresql/sqlc"
	"github.com/KirillZharkov/Ecommerce-API/internal/orders"
	"github.com/KirillZharkov/Ecommerce-API/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
)

type API struct {
	config Config
	db     *pgx.Conn
}

type Config struct {
	addr string //адрес сервера порт 8080
	db   dbConfig
}

type dbConfig struct {
	dsn string //строка домена, для подключения к бд
}

func (app *API) mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second)) //перерыв

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})

	productService := products.NewService(repo.New(app.db))
	productHandler := products.NewHandler(productService)

	router.Get("/products", productHandler.ListProducts)
	router.Get("/products/{id}", productHandler.FindPoductsByID)

	orderServise := orders.NewService(repo.New(app.db), app.db)
	ordersHandler := orders.NewHandler(orderServise)

	router.Post("/orders", ordersHandler.PlaceOrder)

	return router
}

func (app *API) run(h http.Handler) error {

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}
