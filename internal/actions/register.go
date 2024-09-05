package actions

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/spqrcor/gofermart/internal/db"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const minPasswordLength = 6
const maxPasswordLength = 72
const minLoginLength = 3
const maxLoginLength = 255

type InputDataUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var ErrLoginExists = fmt.Errorf("login exists")
var ErrValidation = fmt.Errorf("validation error")

func Register(ctx context.Context, input InputDataUser) (uuid.UUID, error) {
	if err := validate(input); err != nil {
		return uuid.Nil, err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return uuid.Nil, err
	}

	baseUserID := ""
	userId := uuid.NewString()
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	err = db.Source.QueryRowContext(childCtx, "INSERT INTO users (id, login, password) VALUES ($1, $2, $3)  "+
		"ON CONFLICT(login) DO UPDATE SET login = EXCLUDED.login RETURNING id", userId, input.Login, string(bytes)).Scan(&baseUserID)
	if err != nil {
		return uuid.Nil, err
	} else if userId != baseUserID {
		return uuid.Nil, ErrLoginExists
	}
	return uuid.MustParse(baseUserID), nil
}

func validate(input InputDataUser) error {
	if (len(input.Login) < minLoginLength) || (len(input.Login) > maxLoginLength) {
		return fmt.Errorf("%w: ошибка при заполнении login, корректная длина от 3 до 255", ErrValidation)
	}
	if (len(input.Password) < minPasswordLength) || (len(input.Password) > maxPasswordLength) {
		return fmt.Errorf("%w: ошибка при заполнении password, корректная длина от 6 до 72", ErrValidation)
	}
	return nil
}
