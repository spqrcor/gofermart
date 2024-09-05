package actions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/spqrcor/gofermart/internal/db"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrLogin = fmt.Errorf("login error")

func Login(ctx context.Context, input InputDataUser) (uuid.UUID, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	row := db.Source.QueryRowContext(childCtx, "SELECT id, password FROM users WHERE login = $1", input.Login)

	var userID, password string
	err := row.Scan(&userID, &password)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, ErrLogin
	} else if err != nil {
		return uuid.Nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(input.Password)); err != nil {
		return uuid.Nil, ErrLogin
	}
	return uuid.MustParse(userID), nil
}
