package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"purchase-cart-service/dtos"
	"purchase-cart-service/errors"
	"purchase-cart-service/models"
	"strings"
)

type OrderItemRepository struct {
	Db *sql.DB
}

func (bs OrderItemRepository) Insert(ctx context.Context, orderId int, item dtos.OrderItem, product models.Product) (*models.OrderItem, error) {

	query := "INSERT INTO order_item (id, order_id, product_id, price, vat, quantity) VALUES (DEFAULT, $1, $2, $3, $4, $5) RETURNING id, order_id, product_id, price, vat, quantity, created_at"
	slog.Info(query)

	stmt, err := bs.Db.PrepareContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rv := &models.OrderItem{}

	itemPrice := float64(item.Quantity) * product.Price
	itemVAT := itemPrice * (product.VatRate / 100)
	tmp := fmt.Sprint(orderId, item)
	slog.Info(tmp)
	err = stmt.QueryRowContext(ctx, orderId, item.ProductId, itemPrice, itemVAT, item.Quantity).Scan(&rv.Id, &rv.OrderId, &rv.ProductId, &rv.Price, &rv.VAT, &rv.Quantity, &rv.CreatedAt)

	if err != nil {
		return nil, err
	}

	slog.Info("query executed successfully")

	return rv, nil
}

func (bs OrderItemRepository) InsertBatch(ctx context.Context, orderId int, items []dtos.OrderItem, products map[int]models.Product) ([]models.OrderItem, error) {

	query := "INSERT INTO order_item (id, order_id, product_id, price, vat, quantity) VALUES "

	values := []interface{}{}
	placeholders := []string{}

	for i, item := range items {
		product, ok := products[item.ProductId]
		if !ok {
			return nil, errors.PrintAndReturnErr(nil, errors.NewAPIError(http.StatusBadRequest, fmt.Sprintf("product with ID %d not found", item.ProductId)))
		}

		itemPrice := float64(item.Quantity) * product.Price
		itemVAT := itemPrice * (product.VatRate / 100)

		// Create placeholders for this row
		placeholder := fmt.Sprintf("(DEFAULT, $%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)
		placeholders = append(placeholders, placeholder)

		// Add the values to the arguments slice
		values = append(values, orderId, item.ProductId, itemPrice, itemVAT, item.Quantity)
	}

	query += strings.Join(placeholders, ", ")
	query += " RETURNING id, order_id, product_id, price, vat, quantity, created_at"
	slog.Info(query)

	// Execute the query
	rows, err := bs.Db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	slog.Info("query executed successfully")

	insertedItems := []models.OrderItem{}
	for rows.Next() {
		var rv models.OrderItem
		if err := rows.Scan(&rv.Id, &rv.OrderId, &rv.ProductId, &rv.Price, &rv.VAT, &rv.Quantity, &rv.CreatedAt); err != nil {
			return nil, err
		}

		insertedItems = append(insertedItems, rv)
	}

	if err := rows.Err(); err != nil {
		return nil, err

	}

	return insertedItems, nil
}
