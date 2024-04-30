package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// Purchase представляет одну покупку книги пользователем.
type Purchase struct {
	ID         int64     `json:"id"`          // Уникальный идентификатор покупки
	UserID     int64     `json:"user_id"`     // Идентификатор пользователя, совершившего покупку
	BookID     int64     `json:"book_id"`     // Идентификатор купленной книги
	Quantity   int       `json:"quantity"`    // Количество купленных экземпляров
	TotalPrice float64   `json:"total_price"` // Общая цена покупки
	CreatedAt  time.Time `json:"created_at"`  // Время создания записи о покупке
}

type PurchaseModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (m *PurchaseModel) Insert(purchase *Purchase) error {
	query := `
    INSERT INTO purchases (user_id, book_id, quantity, total_price, created_at)
    VALUES ($1, $2, $3, $4, NOW()) RETURNING id;
    `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, purchase.UserID, purchase.BookID, purchase.Quantity, purchase.TotalPrice).Scan(&purchase.ID)
}

func (m PurchaseModel) GetByUserID(userID int64) ([]*Book, error) {
	query := `
	SELECT b.id, b.title, b.author, b.price, b.stock_quantity, b.created_at, b.updated_at
	FROM purchases p
	INNER JOIN books b ON p.book_id = b.id
	WHERE p.user_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Price,
			&book.StockQuantity,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (m *PurchaseModel) GetAllForUser(userID int64) ([]*Purchase, error) {
	query := `
    SELECT id, user_id, book_id, quantity, total_price, created_at
    FROM purchases
    WHERE user_id = $1;
    `
	var purchases []*Purchase

	rows, err := m.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Purchase
		err := rows.Scan(&p.ID, &p.UserID, &p.BookID, &p.Quantity, &p.TotalPrice, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return purchases, nil
}
