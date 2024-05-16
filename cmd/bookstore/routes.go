package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {
	r := mux.NewRouter()
	// Convert the app.notFoundResponse helper to a http.Handler using the http.HandlerFunc()
	// adapter, and then set it as the custom error handler for 404 Not Found responses.
	r.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	// Convert app.methodNotAllowedResponse helper to a http.Handler and set it as the custom
	// error handler for 405 Method Not Allowed responses
	r.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	// Настройка маршрутов для книг
	bookRouter := r.PathPrefix("/api/v1/books").Subrouter()
	bookRouter.HandleFunc("", app.createBookHandler).Methods("POST")               // Создание книги
	bookRouter.HandleFunc("/{id:[0-9]+}", app.getBookHandler).Methods("GET")       // Получение книги по ID
	bookRouter.HandleFunc("/{id:[0-9]+}", app.updateBookHandler).Methods("PUT")    // Обновление книги
	bookRouter.HandleFunc("/{id:[0-9]+}", app.deleteBookHandler).Methods("DELETE") // Удаление книги
	bookRouter.HandleFunc("/buy", app.requirePermissions("books:read", app.BuyBook)).Methods("POST")
	bookRouter.HandleFunc("/list", app.GetBookList).Methods("GET")

	//Users handlers
	users1 := r.PathPrefix("/api/v1/users").Subrouter()
	// User handlers with Authentication
	users1.HandleFunc("", app.registerUserHandler).Methods("POST")
	users1.HandleFunc("/activated", app.activateUserHandler).Methods("PUT")
	users1.HandleFunc("/login", app.createAuthenticationTokenHandler).Methods("POST")
	users1.HandleFunc("/purchases", app.requirePermissions("books:read", app.ListPurchases)).Methods("GET")

	commentsRouter := r.PathPrefix("/api/v1/comments").Subrouter()
	commentsRouter.HandleFunc("", app.CreateComment).Methods("POST")
	commentsRouter.HandleFunc("", app.GetComments).Methods("GET")
	commentsRouter.HandleFunc("", app.DeleteComment).Methods("DELETE")

	// Wrap the router with the panic recovery middleware and rate limit middleware.
	return app.authenticate(r)
}
