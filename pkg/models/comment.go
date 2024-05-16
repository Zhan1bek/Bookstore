package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
)

type Comment struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	BookID    int64     `json:"book_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// Insert добавляет новый комментарий в базу данных.
func (m *CommentModel) Insert(comment *Comment) error {
	query := `
		INSERT INTO comments (user_id, book_id, content) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`
	args := []interface{}{comment.UserID, comment.BookID, comment.Content}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		m.ErrorLog.Printf("Ошибка при вставке нового комментария: %v", err)
		return err
	}
	m.InfoLog.Printf("Комментарий успешно добавлен с ID %d", comment.ID)
	return nil
}

// Get получает комментарий по его ID.
func (m *CommentModel) Get(id int64) (*Comment, error) {
	query := `
		SELECT id, user_id, book_id, content, created_at
		FROM comments
		WHERE id = $1
	`
	var comment Comment
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&comment.ID, &comment.UserID, &comment.BookID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetAllByBook получает все комментарии к заданной книге.
func (m *CommentModel) GetAllByBook(bookID int64) ([]*Comment, error) {
	query := `
		SELECT id, user_id, book_id, content, created_at
		FROM comments
		WHERE book_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.BookID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// Update обновляет содержание комментария.
func (m *CommentModel) Update(comment *Comment) error {
	query := `
		UPDATE comments
		SET content = $1
		WHERE id = $2
		RETURNING created_at
	`
	args := []interface{}{comment.Content, comment.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&comment.CreatedAt)
}

// Delete удаляет комментарий по его ID.
// Delete удаляет комментарий по его ID, если он принадлежит указанному пользователю.
func (m *CommentModel) Delete(commentID, userID int64) error {
	query := `
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, commentID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no matching record found or user not authorized")
	}

	return nil
}
