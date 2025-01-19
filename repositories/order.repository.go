package repositories

import (
	"context"
	"database/sql"
	"log/slog"
	"purchase-cart-service/errors"
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
		return nil, errors.PrintAndReturnErr(err, errors.INTERNAL_SERVER_ERROR)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, time.Now()).Scan(orderId)
	if err != nil {
		return nil, errors.PrintAndReturnErr(err, errors.INTERNAL_SERVER_ERROR)
	}

	slog.Info("record inserted successfully")
	return orderId, nil
}
