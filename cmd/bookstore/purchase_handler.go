package main

import (
	"encoding/json"
	"errors"
	"github.com/Zhan1bek/BookStore/pkg/models"
	"net/http"
)

func (app *application) BuyBook(w http.ResponseWriter, r *http.Request) {
	token, err := app.GetToken(w, r)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	user, err := app.models.Users.GetByToken(models.ScopeAuthentication, token)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.invalidAuthenticationTokenResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title string `json:"title"` // Assume JSON body contains a 'title' field
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	book, err := app.models.Books.GetByTitle(input.Title)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if book.StockQuantity < 1 {
		app.errorResponse(w, r, http.StatusForbidden, "No more books available in stock.")
		return
	}

	book.StockQuantity -= 1
	err = app.models.Books.UpdateQuantity(book.ID, book.StockQuantity)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	purchase := &models.Purchase{
		UserID:     user.ID,
		BookID:     book.ID,
		Quantity:   1,
		TotalPrice: book.Price,
	}
	err = app.models.Purchase.Insert(purchase)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"purchase": purchase}, nil)
}
