package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"purchase-cart-service/models"
	"time"

	"github.com/lib/pq"
)

type ProductRepository struct {
	Db *sql.DB
}

func (bs ProductRepository) Insert(ctx context.Context, product models.Product) (*models.Product, error) {

	stmt, err := bs.Db.PrepareContext(ctx, "INSERT INTO product (id, name, price, vat_rate, created_at) VALUES (DEFAULT, $1, $2, $3, $4)")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, product.Name, product.Price, product.VatRate, time.Now())

	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (bs ProductRepository) GetByID(ctx context.Context, id int) (*models.Product, error) {
	// Prepare query with a placeholder for the IN clause
	query := "SELECT id, name, description, price, vat_rate, created_at FROM product WHERE id = $1"

	// Execute the query
	rows, err := bs.Db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Initialize a slice to store the results
	var product *models.Product

	// Iterate over the rows and map them to the Product model
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.VatRate,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return product, nil
}

func (bs ProductRepository) GetByIDs(ctx context.Context, ids []int) ([]models.Product, error) {
	query := `SELECT id, "name", price, vat_rate, created_at FROM product where id = ANY ($1)`

	slog.Info(query)

	rows, err := bs.Db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		slog.Error("Failed to execute query")
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.VatRate,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(products) != len(ids) {
		return nil, errors.New("invalid products in items list")
	}

	slog.Info("query executed successfully")
	return products, nil
}

// func (bs ProductRepository) GetAll(ctx context.Context) ([]models.Book, error) {

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
