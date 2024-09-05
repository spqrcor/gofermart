package actions

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/db"
	"github.com/spqrcor/gofermart/internal/logger"
	"time"
)

type Order struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual,omitempty"`
	UploadedAt string `json:"uploaded_at"`
}

var ErrOrdersNotFound = fmt.Errorf("orders not found")

func GetOrders(ctx context.Context) ([]Order, error) {
	var orders []Order
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	rows, err := db.Source.QueryContext(childCtx, "SELECT number, status, accrual, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC", ctx.Value(authenticate.ContextUserID))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Log.Error(err.Error())
		}
		if err := rows.Err(); err != nil {
			logger.Log.Error(err.Error())
		}
	}()

	for rows.Next() {
		o := Order{}
		var accrual sql.NullInt32
		if err = rows.Scan(&o.Number, &o.Status, &accrual, &o.UploadedAt); err != nil {
			return nil, err
		}
		o.Accrual = int(accrual.Int32)
		orders = append(orders, o)
	}

	if len(orders) == 0 {
		return nil, ErrOrdersNotFound
	}
	return orders, nil
}
