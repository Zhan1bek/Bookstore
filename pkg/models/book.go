package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BookModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *BookModel) Insert(book *Book) error {
	query := `
        INSERT INTO books (title, author, price, stock_quantity) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id, created_at, updated_at
    `
	args := []interface{}{book.Title, book.Author, book.Price, book.StockQuantity}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
	if err != nil {
		m.ErrorLog.Printf("Ошибка при вставке новой книги: %v", err)
		return err
	}
	m.InfoLog.Printf("Книга [%s] успешно добавлена с ID %d", book.Title, book.ID)
	return nil
}

func (m BookModel) Get(id int) (*Book, error) {
	query := `
    SELECT id, created_at, updated_at, title, author, price, stock_quantity
    FROM books
    WHERE id = $1
    `
	var book Book
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt, &book.Title, &book.Author, &book.Price, &book.StockQuantity)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (m BookModel) Update(book *Book) error {
	query := `
    UPDATE books
    SET title = $1, author = $2, price = $3, stock_quantity = $4
    WHERE id = $5
    RETURNING updated_at
    `
	args := []interface{}{book.Title, book.Author, book.Price, book.StockQuantity, book.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&book.UpdatedAt)
}

func (m BookModel) Delete(id int) error {
	query := `
    DELETE FROM books
    WHERE id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}
