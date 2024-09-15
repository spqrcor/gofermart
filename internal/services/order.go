package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/db"
	"github.com/spqrcor/gofermart/internal/logger"
	"github.com/spqrcor/gofermart/internal/utils"
	"time"
)

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

type OrderFromAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

var ErrOrderAnotherUserExists = fmt.Errorf("order another user exists")
var ErrOrderUserExists = fmt.Errorf("order user exists")
var ErrOrderInvalidFormat = fmt.Errorf("order invalid format")
var ErrOrdersNotFound = fmt.Errorf("orders not found")

type OrderRepository interface {
	Add(ctx context.Context, orderNum string) error
	GetAll(ctx context.Context) ([]Order, error)
	GetUnComplete(ctx context.Context) ([]string, error)
	ChangeStatus(ctx context.Context, data OrderFromAccrual) error
}

type OrderService struct{}

func NewOrderService() *OrderService {
	return &OrderService{}
}

func (o *OrderService) Add(ctx context.Context, orderNum string) error {
	if !utils.IsNumberValid(orderNum) {
		return ErrOrderInvalidFormat
	}
	var baseUserID, baseOrderID string
	orderID := uuid.NewString()

	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	err := db.Source.QueryRowContext(childCtx, "INSERT INTO orders (id, user_id, number) VALUES ($1, $2, $3)  "+
		"ON CONFLICT(number) DO UPDATE SET number = EXCLUDED.number RETURNING id, user_id", orderID, ctx.Value(authenticate.ContextUserID), orderNum).Scan(&baseOrderID, &baseUserID)
	if err != nil {
		return err
	} else if ctx.Value(authenticate.ContextUserID) != uuid.MustParse(baseUserID) {
		return ErrOrderAnotherUserExists
	} else if orderID != baseOrderID {
		return ErrOrderUserExists
	}
	return nil
}

func (o *OrderService) GetAll(ctx context.Context) ([]Order, error) {
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
		var accrual sql.NullFloat64
		if err = rows.Scan(&o.Number, &o.Status, &accrual, &o.UploadedAt); err != nil {
			return nil, err
		}
		o.Accrual = accrual.Float64
		orders = append(orders, o)
	}

	if len(orders) == 0 {
		return nil, ErrOrdersNotFound
	}
	return orders, nil
}

func (o *OrderService) GetUnComplete(ctx context.Context) ([]string, error) {
	var orders []string

	childCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := db.Source.QueryContext(childCtx, "SELECT number FROM orders WHERE status IN ('NEW', 'PROCESSING') ORDER BY created_at")
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
		var orderNum string
		if err = rows.Scan(&orderNum); err != nil {
			return nil, err
		}
		orders = append(orders, orderNum)
	}
	return orders, nil
}

func (o *OrderService) ChangeStatus(ctx context.Context, data OrderFromAccrual) error {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	tx, err := db.Source.BeginTx(childCtx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(childCtx, "UPDATE orders SET status = $1, accrual =$2 WHERE number = $3", data.Status, data.Accrual, data.Order)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if data.Accrual > 0 {
		_, err = tx.ExecContext(childCtx, "UPDATE users SET balance = balance + $1 WHERE id = (SELECT user_id FROM orders WHERE number = $2)", data.Accrual, data.Order)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
