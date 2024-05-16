package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Book представляет модель книги
type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	AvgRating     float64   `json:"avg_rating"`
	RatingCount   int       `json:"rating_count"`
}

// BookModel обрабатывает операции с книгами в базе данных
type BookModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// GetAll возвращает все книги с учетом фильтров
func (m BookModel) GetAll(title string, author string, priceRangeStart, priceRangeEnd, minRating float64, filters Filters) ([]*Book, Metadata, error) {
	// Конструирование SQL-запроса с фильтрацией и сортировкой
	query := fmt.Sprintf(`
        SELECT count(*) OVER(), id, created_at, updated_at, title, author, price, stock_quantity, avg_rating, rating_count
        FROM books
        WHERE (LOWER(title) LIKE LOWER($1) OR $1 = '')
        AND (LOWER(author) LIKE LOWER($2) OR $2 = '')
        AND (price >= $3 OR $3 = 0)
        AND (price <= $4 OR $4 = 0)
        AND (avg_rating >= $5 OR $5 = 0)
        ORDER BY %s %s, id ASC
        LIMIT $6 OFFSET $7
        `, filters.sortColumn(), filters.sortDirection())

	// Создание контекста с тайм-аутом 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Подготовка аргументов для запроса
	args := []interface{}{"%" + title + "%", "%" + author + "%", priceRangeStart, priceRangeEnd, minRating, filters.limit(), filters.offset()}

	// Выполнение запроса
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			m.ErrorLog.Println(err)
		}
	}()

	// Инициализация переменных для сканирования строк
	totalRecords := 0
	var books []*Book

	// Итерация по результирующему набору и сканирование каждой строки в структуру Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&totalRecords, &book.ID, &book.CreatedAt, &book.UpdatedAt, &book.Title, &book.Author, &book.Price, &book.StockQuantity, &book.AvgRating, &book.RatingCount)
		if err != nil {
			return nil, Metadata{}, err
		}
		books = append(books, &book)
	}

	// Проверка на наличие ошибок во время итерации
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Подготовка метаданных
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return books, metadata, nil
}

// Insert вставляет новую книгу в базу данных
func (m *BookModel) Insert(book *Book) error {
	query := `
        INSERT INTO books (title, author, price, stock_quantity) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id, created_at, updated_at, avg_rating, rating_count
    `
	args := []interface{}{book.Title, book.Author, book.Price, book.StockQuantity}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt, &book.AvgRating, &book.RatingCount)
	if err != nil {
		m.ErrorLog.Printf("Ошибка при вставке новой книги: %v", err)
		return err
	}
	m.InfoLog.Printf("Книга [%s] успешно добавлена с ID %d", book.Title, book.ID)
	return nil
}

// Get возвращает книгу по ID
func (m BookModel) Get(id int64) (*Book, error) {
	query := `
    SELECT id, created_at, updated_at, title, author, price, stock_quantity, avg_rating, rating_count
    FROM books
    WHERE id = $1
    `
	var book Book
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt, &book.Title, &book.Author, &book.Price, &book.StockQuantity, &book.AvgRating, &book.RatingCount)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Update обновляет информацию о книге в базе данных
func (m BookModel) Update(book *Book) error {
	query := `
    UPDATE books
    SET title = $1, author = $2, price = $3, stock_quantity = $4
    WHERE id = $5
    RETURNING updated_at, avg_rating, rating_count
    `
	args := []interface{}{book.Title, book.Author, book.Price, book.StockQuantity, book.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&book.UpdatedAt, &book.AvgRating, &book.RatingCount)
}

// Delete удаляет книгу из базы данных
func (m BookModel) Delete(id int64) error {
	query := `
    DELETE FROM books
    WHERE id = $1
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// GetByTitle возвращает книгу по названию
func (m *BookModel) GetByTitle(title string) (*Book, error) {
	query := `
    SELECT id, title, author, price, stock_quantity, created_at, updated_at, avg_rating, rating_count
    FROM books
    WHERE title = $1
    `
	var book Book
	err := m.DB.QueryRow(query, title).Scan(
		&book.ID, &book.Title, &book.Author, &book.Price, &book.StockQuantity, &book.CreatedAt, &book.UpdatedAt, &book.AvgRating, &book.RatingCount,
	)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// UpdateRating обновляет средний рейтинг и количество оценок книги
func (m BookModel) UpdateRating(bookID int64) error {
	query := `
        UPDATE books
        SET avg_rating = (
            SELECT AVG(rating) FROM ratings WHERE book_id = $1
        ),
        rating_count = (
            SELECT COUNT(*) FROM ratings WHERE book_id = $1
        )
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, bookID)
	return err
}

// UpdateQuantity обновляет количество книг на складе
func (m *BookModel) UpdateQuantity(bookID int64, newQuantity int) error {
	query := `
    UPDATE books
    SET stock_quantity = $1
    WHERE id = $2;
    `
	_, err := m.DB.Exec(query, newQuantity, bookID)
	return err
}
