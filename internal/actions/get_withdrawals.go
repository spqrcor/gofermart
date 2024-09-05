package actions

import (
	"context"
	"fmt"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/db"
	"github.com/spqrcor/gofermart/internal/logger"
	"time"
)

type Withdraw struct {
	OrderNum    string `json:"order"`
	Sum         int    `json:"sum"`
	ProcessedAt string `json:"processed_at"`
}

var ErrWithdrawNotFound = fmt.Errorf("withdraw not found")

func GetWithdraws(ctx context.Context) ([]Withdraw, error) {
	var withdraws []Withdraw
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	rows, err := db.Source.QueryContext(childCtx, "SELECT o.number, wl.sum, wl.created_at FROM withdraw_list wl "+
		"INNER JOIN orders o ON o.id = wl.order_id WHERE o.user_id = $1 ORDER BY wl.created_at DESC", ctx.Value(authenticate.ContextUserID))
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
		w := Withdraw{}
		if err = rows.Scan(&w.OrderNum, &w.Sum, &w.ProcessedAt); err != nil {
			return nil, err
		}
		withdraws = append(withdraws, w)
	}
	if len(withdraws) == 0 {
		return nil, ErrWithdrawNotFound
	}
	return withdraws, nil
}
