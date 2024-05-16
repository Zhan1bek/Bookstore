package models

import (
	"database/sql"
	"log"
	"time"
)

// Rating представляет рейтинг книги от пользователя
type Rating struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	BookID    int       `json:"book_id"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

type RatingModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// Insert вставляет новый рейтинг в базу данных
func (m RatingModel) Insert(rating *Rating) error {
	query := `
        INSERT INTO ratings (user_id, book_id, rating)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`
	args := []interface{}{rating.UserID, rating.BookID, rating.Rating}

	return m.DB.QueryRow(query, args...).Scan(&rating.ID, &rating.CreatedAt)
}
