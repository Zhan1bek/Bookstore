package main

import (
	"encoding/json"
	"github.com/Zhan1bek/BookStore/pkg/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title         string  `json:"title"`
		Author        string  `json:"author"`
		Price         float64 `json:"price"`
		StockQuantity int     `json:"stock_quantity"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book := &models.Book{
		Title:         input.Title,
		Author:        input.Author,
		Price:         input.Price,
		StockQuantity: input.StockQuantity,
	}

	err = app.models.Books.Insert(book)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, book)
}

// Обработчик для получения книги по ID
func (app *application) getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book, err := app.models.Books.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Book not found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, book)
}

// Обработчик для обновления книги
func (app *application) updateBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	var input struct {
		Title         *string  `json:"title"`
		Author        *string  `json:"author"`
		Price         *float64 `json:"price"`
		StockQuantity *int     `json:"stock_quantity"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	book := models.Book{ID: int64(id)}
	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.Author != nil {
		book.Author = *input.Author
	}
	if input.Price != nil {
		book.Price = *input.Price
	}
	if input.StockQuantity != nil {
		book.StockQuantity = *input.StockQuantity
	}

	err = app.models.Books.Update(&book)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, book)
}

// Обработчик для удаления книги
func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	err = app.models.Books.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
