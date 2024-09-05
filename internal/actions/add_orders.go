package actions

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/spqrcor/gofermart/internal/authenticate"
	"github.com/spqrcor/gofermart/internal/db"
	"time"
)

var ErrOrderAnotherUserExists = fmt.Errorf("order another user exists")
var ErrOrderUserExists = fmt.Errorf("order user exists")
var ErrOrderInvalidFormat = fmt.Errorf("order invalid format")

func AddOrder(ctx context.Context, orderNum string) error {
	if !isNumberValid(orderNum) {
		return ErrOrderInvalidFormat
	}

	var baseUserID, baseOrderID string
	orderId := uuid.NewString()
	userId := ctx.Value(authenticate.ContextUserID)
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	err := db.Source.QueryRowContext(childCtx, "INSERT INTO orders (id, user_id, number) VALUES ($1, $2, $3)  "+
		"ON CONFLICT(number) DO UPDATE SET number = EXCLUDED.number RETURNING id, user_id", orderId, userId, orderNum).Scan(&baseOrderID, &baseUserID)
	if err != nil {
		return err
	} else if userId != uuid.MustParse(baseUserID) {
		return ErrOrderAnotherUserExists
	} else if orderId != baseOrderID {
		return ErrOrderUserExists
	}
	return nil
}

func isNumberValid(orderNum string) bool {
	total := 0
	isSecondDigit := false
	for i := len(orderNum) - 1; i >= 0; i-- {
		digit := int(orderNum[i] - '0')
		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		total += digit
		isSecondDigit = !isSecondDigit
	}
	return total%10 == 0
}
