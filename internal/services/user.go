package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrLogin = fmt.Errorf("login error")
var ErrLoginExists = fmt.Errorf("login exists")
var ErrValidation = fmt.Errorf("validation error")
var ErrGeneratePassword = fmt.Errorf("generate password error")

const minPasswordLength = 6
const maxPasswordLength = 72
const minLoginLength = 3
const maxLoginLength = 255

type InputDataUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserService struct {
	db           *sql.DB
	queryTimeOut time.Duration
}

func NewUserService(db *sql.DB, queryTimeOut time.Duration) *UserService {
	return &UserService{db: db, queryTimeOut: queryTimeOut}
}

func (u *UserService) Add(ctx context.Context, input InputDataUser) (uuid.UUID, error) {
	if err := validate(input); err != nil {
		return uuid.Nil, err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return uuid.Nil, errors.Join(ErrGeneratePassword, err)
	}

	baseUserID := ""
	userID := uuid.NewString()
	childCtx, cancel := context.WithTimeout(ctx, time.Second*u.queryTimeOut)
	defer cancel()
	err = u.db.QueryRowContext(childCtx, "INSERT INTO users (id, login, password) VALUES ($1, $2, $3) ON CONFLICT(login) DO UPDATE SET login = EXCLUDED.login RETURNING id",
		userID, input.Login, string(bytes)).Scan(&baseUserID)
	if err != nil {
		return uuid.Nil, err
	}
	if userID != baseUserID {
		return uuid.Nil, ErrLoginExists
	}
	return uuid.MustParse(baseUserID), nil
}

func (u *UserService) Login(ctx context.Context, input InputDataUser) (uuid.UUID, error) {
	childCtx, cancel := context.WithTimeout(ctx, time.Second*u.queryTimeOut)
	defer cancel()
	row := u.db.QueryRowContext(childCtx, "SELECT id, password FROM users WHERE login = $1", input.Login)

	var userID, password string
	err := row.Scan(&userID, &password)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, ErrLogin
	}
	if err != nil {
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
