package repositories

import (
	"context"
	"database/sql"
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
	query := "SELECT id, name, description, price, vat_rate, created_at FROM product WHERE id = $1"

	rows, err := bs.Db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err

	}
	defer rows.Close()

	var product *models.Product

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

	slog.Info("query executed successfully")
	return products, nil
}
