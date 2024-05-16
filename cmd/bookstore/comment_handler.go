package main

import (
	"encoding/json"
	"errors"
	"github.com/Zhan1bek/BookStore/pkg/models"
	"net/http"
	"strconv"
)

// CreateComment создает новый комментарий к книге.
func (app *application) CreateComment(w http.ResponseWriter, r *http.Request) {
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
		BookID  int64  `json:"book_id"`
		Content string `json:"content"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &models.Comment{
		UserID:  user.ID,
		BookID:  input.BookID,
		Content: input.Content,
	}
	err = app.models.Comment.Insert(comment)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"comment": comment}, nil)
}

// GetComments получает все комментарии к книге.
func (app *application) GetComments(w http.ResponseWriter, r *http.Request) {
	token, err := app.GetToken(w, r)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	_, err = app.models.Users.GetByToken(models.ScopeAuthentication, token)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.invalidAuthenticationTokenResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	bookIDStr := r.URL.Query().Get("book_id")
	if bookIDStr == "" {
		app.errorResponse(w, r, http.StatusBadRequest, "Book ID is required")
		return
	}

	bookID, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid Book ID")
		return
	}

	comments, err := app.models.Comment.GetAllByBook(bookID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"comments": comments}, nil)
}

// DeleteComment удаляет комментарий, если пользователь является его автором.
func (app *application) DeleteComment(w http.ResponseWriter, r *http.Request) {
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
		CommentID int64 `json:"comment_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Comment.Delete(input.CommentID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, errors.New("no matching record found or user not authorized")):
			app.errorResponse(w, r, http.StatusForbidden, "You do not have permission to delete this comment.")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
