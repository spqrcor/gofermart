package actions

import (
	"context"
	"errors"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/db"
	"time"
)

type BalanceInfo struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

func GetBalanceInfo(ctx context.Context) (BalanceInfo, error) {
	balanceInfo := BalanceInfo{}
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	row := db.Source.QueryRowContext(childCtx, "SELECT balance, (SELECT COALESCE(SUM(wl.sum), 0) FROM withdraw_list wl INNER JOIN orders o ON o.id = wl.order_id WHERE o.user_id = u.id) "+
		"FROM users u WHERE id = $1", ctx.Value(authenticate.ContextUserID))
	if err := row.Scan(&balanceInfo.Current, &balanceInfo.Withdrawn); err != nil {
		return balanceInfo, errors.New("User not found")
	}
	return balanceInfo, nil
}
