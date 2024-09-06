package actions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/db"
	"time"
)

type InputWithdraw struct {
	OrderNum string `json:"order"`
	Sum      int    `json:"sum"`
}

var ErrBalance = fmt.Errorf("balance error")

func AddWithdraw(ctx context.Context, input InputWithdraw) error {
	if !isNumberValid(input.OrderNum) {
		return ErrOrderInvalidFormat
	}
	userID := ctx.Value(authenticate.ContextUserID)

	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	tx, err := db.Source.BeginTx(childCtx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(childCtx, "UPDATE users SET balance = balance - $2 WHERE id = $1", userID, input.Sum)
	if err != nil {
		_ = tx.Rollback()
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.ConstraintName == "users_balance_check" {
			return ErrBalance
		}
		return err
	}

	var withdrawID string
	err = tx.QueryRowContext(childCtx, "INSERT INTO withdraw_list (order_id, sum) SELECT id, $3 FROM orders WHERE user_id = $1 AND number = $2 RETURNING id", userID, input.OrderNum, input.Sum).Scan(&withdrawID)
	if err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return ErrOrderInvalidFormat
		}
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
