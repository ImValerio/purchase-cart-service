package repositories

import (
	"context"
	"database/sql"
	"log/slog"
	"time"
)

type OrderRepository struct {
	Db *sql.DB
}

func (bs OrderRepository) Insert(ctx context.Context) (*int, error) {

	query := `INSERT INTO "order" (id, created_at) VALUES (DEFAULT, $1) RETURNING id`

	slog.Info(query)
	orderId := new(int)

	stmt, err := bs.Db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, time.Now()).Scan(orderId)
	if err != nil {
		return nil, err
	}

	slog.Info("record inserted successfully")
	return orderId, nil
}

// func (bs OrderRepository) GetAll(ctx context.Context) ([]models.Book, error) {

// 	stmt, err := bs.Db.PrepareContext(ctx, "SELECT id, name, author, create_time FROM book")

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Query()

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var books []models.Book

// 	for rows.Next() {
// 		var book models.Book

// 		err := rows.Scan(&book.Id, &book.Name, &book.Author, &book.CreateTime)

// 		if err != nil {
// 			return nil, err
// 		}

// 		books = append(books, book)
// 	}

// 	return books, nil
// }
