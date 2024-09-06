package services

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
var ErrLoginExists = fmt.Errorf("login exists")
var ErrValidation = fmt.Errorf("validation error")

const minPasswordLength = 6
const maxPasswordLength = 72
const minLoginLength = 3
const maxLoginLength = 255

type InputDataUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserRepository interface {
	Add(ctx context.Context, input InputDataUser) (uuid.UUID, error)
	Login(ctx context.Context, input InputDataUser) (uuid.UUID, error)
}

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) Add(ctx context.Context, input InputDataUser) (uuid.UUID, error) {
	if err := validate(input); err != nil {
		return uuid.Nil, err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return uuid.Nil, err
	}

	baseUserID := ""
	userID := uuid.NewString()
	childCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	err = db.Source.QueryRowContext(childCtx, "INSERT INTO users (id, login, password) VALUES ($1, $2, $3) ON CONFLICT(login) DO UPDATE SET login = EXCLUDED.login RETURNING id",
		userID, input.Login, string(bytes)).Scan(&baseUserID)
	if err != nil {
		return uuid.Nil, err
	} else if userID != baseUserID {
		return uuid.Nil, ErrLoginExists
	}
	return uuid.MustParse(baseUserID), nil
}

func (u *UserService) Login(ctx context.Context, input InputDataUser) (uuid.UUID, error) {
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

func validate(input InputDataUser) error {
	if (len(input.Login) < minLoginLength) || (len(input.Login) > maxLoginLength) {
		return fmt.Errorf("%w: ошибка при заполнении login, корректная длина от 3 до 255", ErrValidation)
	}
	if (len(input.Password) < minPasswordLength) || (len(input.Password) > maxPasswordLength) {
		return fmt.Errorf("%w: ошибка при заполнении password, корректная длина от 6 до 72", ErrValidation)
	}
	return nil
}
