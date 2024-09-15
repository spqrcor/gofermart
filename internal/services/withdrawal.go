package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/db"
	"github.com/spqrcor/gofermart/internal/logger"
	"github.com/spqrcor/gofermart/internal/utils"
	"time"
)

type InputWithdrawal struct {
	OrderNum string `json:"order"`
	Sum      int    `json:"sum"`
}

type Withdrawal struct {
	OrderNum    string `json:"order"`
	Sum         int    `json:"sum"`
	ProcessedAt string `json:"processed_at"`
}

type BalanceInfo struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

var ErrBalance = fmt.Errorf("balance error")
var ErrWithdrawNotFound = fmt.Errorf("withdraw not found")

type WithdrawalRepository interface {
	Add(ctx context.Context, input InputWithdrawal) error
	GetAll(ctx context.Context) ([]Withdrawal, error)
	GetBalance(ctx context.Context) (BalanceInfo, error)
}

type WithdrawalService struct{}

func NewWithdrawalService() *WithdrawalService {
	return &WithdrawalService{}
}

func (w *WithdrawalService) Add(ctx context.Context, input InputWithdrawal) error {
	if !utils.IsNumberValid(input.OrderNum) {
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
	err = tx.QueryRowContext(childCtx, "INSERT INTO withdrawals (user_id, number, sum) VALUES ($1, $2, $3) RETURNING id", userID, input.OrderNum, input.Sum).Scan(&withdrawID)
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

func (w *WithdrawalService) GetAll(ctx context.Context) ([]Withdrawal, error) {
	var withdrawals []Withdrawal
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	rows, err := db.Source.QueryContext(childCtx, "SELECT number, sum, created_at FROM withdrawals WHERE user_id = $1 ORDER BY created_at DESC",
		ctx.Value(authenticate.ContextUserID))
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
		w := Withdrawal{}
		if err = rows.Scan(&w.OrderNum, &w.Sum, &w.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, w)
	}
	if len(withdrawals) == 0 {
		return nil, ErrWithdrawNotFound
	}
	return withdrawals, nil
}

func (w *WithdrawalService) GetBalance(ctx context.Context) (BalanceInfo, error) {
	balanceInfo := BalanceInfo{}
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	row := db.Source.QueryRowContext(childCtx, "SELECT balance, (SELECT COALESCE(SUM(w.sum), 0) FROM withdrawals w WHERE w.user_id = u.id) FROM users u WHERE id = $1",
		ctx.Value(authenticate.ContextUserID))
	if err := row.Scan(&balanceInfo.Current, &balanceInfo.Withdrawn); err != nil {
		return balanceInfo, errors.New("user not found")
	}
	return balanceInfo, nil
}
