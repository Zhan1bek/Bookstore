package main

import (
	"database/sql"
	"flag"
	"github.com/Zhan1bek/BookStore/pkg/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models models.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:3711157al@localhost/bookstore?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: models.NewModels(db),
	}

	app.run()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *application) run() {
	r := mux.NewRouter()

	// Настройка маршрутов для книг
	bookRouter := r.PathPrefix("/api/v1/books").Subrouter()
	bookRouter.HandleFunc("", app.createBookHandler).Methods("POST")               // Создание книги
	bookRouter.HandleFunc("/{id:[0-9]+}", app.getBookHandler).Methods("GET")       // Получение книги по ID
	bookRouter.HandleFunc("/{id:[0-9]+}", app.updateBookHandler).Methods("PUT")    // Обновление книги
	bookRouter.HandleFunc("/{id:[0-9]+}", app.deleteBookHandler).Methods("DELETE") // Удаление книги

	// Запуск сервера
	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	if err != nil {
		log.Fatal(err)
	}
}
