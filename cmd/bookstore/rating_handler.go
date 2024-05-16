package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Zhan1bek/BookStore/pkg/models"
)

func (app *application) rateBook(w http.ResponseWriter, r *http.Request) {
	// Извлечение токена из запроса
	token, err := app.GetToken(w, r)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// Получение пользователя по токену
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

	// Декодирование тела запроса
	var input struct {
		BookID int `json:"book_id"` // JSON тело содержит поле 'book_id'
		Rating int `json:"rating"`  // JSON тело содержит поле 'rating'
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Проверка валидности рейтинга
	if input.Rating < 1 || input.Rating > 5 {
		app.errorResponse(w, r, http.StatusBadRequest, "Rating must be between 1 and 5.")
		return
	}

	// Получение книги по ID
	book, err := app.models.Books.Get(int64(input.BookID))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Создание нового рейтинга
	rating := &models.Rating{
		UserID: int(user.ID),
		BookID: int(book.ID),
		Rating: input.Rating,
	}
	err = app.models.Rating.Insert(rating)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Обновление среднего рейтинга и количества оценок книги
	err = app.models.Books.UpdateRating(int64(int(book.ID)))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Ответ с успешным добавлением рейтинга
	app.writeJSON(w, http.StatusCreated, envelope{"rating": rating}, nil)
}
