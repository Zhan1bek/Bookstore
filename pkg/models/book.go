package models

import (
	"context"
	"database/sql"
	"fmt"
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

func (m BookModel) GetAll(title string, author string, priceRangeStart, priceRangeEnd float64,
	filters Filters) ([]*Book, Metadata, error) {

	// Construct the SQL query with filtering and ordering.
	query := fmt.Sprintf(
		`
        SELECT count(*) OVER(), id, created_at, updated_at, title, author, price, stock_quantity
        FROM books
        WHERE (LOWER(title) LIKE LOWER($1) OR $1 = '')
        AND (LOWER(author) LIKE LOWER($2) OR $2 = '')
        AND (price >= $3 OR $3 = 0)
        AND (price <= $4 OR $4 = 0)
        ORDER BY %s %s, id ASC
        LIMIT $5 OFFSET $6
        `,
		filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the arguments for the query.
	args := []interface{}{"%" + title + "%", "%" + author + "%", priceRangeStart, priceRangeEnd, filters.limit(), filters.offset()}

	// Execute the query.
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	// Initialize the variables for scanning the rows.
	totalRecords := 0
	var books []*Book

	// Iterate over the result set and scan each row into a Book struct.
	for rows.Next() {
		var book Book
		err := rows.Scan(&totalRecords, &book.ID, &book.CreatedAt, &book.UpdatedAt, &book.Title, &book.Author, &book.Price, &book.StockQuantity)
		if err != nil {
			return nil, Metadata{}, err
		}
		books = append(books, &book)
	}

	// Check for any error that occurred during iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Prepare the metadata.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return books, metadata, nil
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
