package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zrotrasukha/jobman/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int32     `json:"_"`
}

type Password struct {
	plainText *string
	hash      []byte
}

func (p *Password) Set(plainTextPassword string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = b

	return nil
}

func (p *Password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidatePlainText(v *validator.Validator, Password string) {
	v.CheckField(Password != "", "password", "must be provided")
	v.CheckField(len(Password) >= 8, "password", "must be at least 8 bytes long")
	v.CheckField(len(Password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.CheckField(email != "", "email", "must be provided")
	v.CheckField(len(email) <= 500, "email", "must not be more than 500 bytes long")
	v.CheckField(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.CheckField(user.Name != "", "name", "must be provided")
	v.CheckField(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plainText != nil {
		ValidatePlainText(v, *user.Password.plainText)
	}

	if user.Password.hash != nil {
		panic("missing Password hash for user")
	}
}

type UserModelInterface interface {
	Insert(user *User) error
}

type UserModel struct {
	pool *pgxpool.Pool
}

func (m UserModel) Insert(user *User) error {
	query := `INSERT INTO users (name, email, Password_hash, activated)
						VALUES ($1, $2, $3, $4)
						RETURNING id, created_at, version`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.pool.QueryRow(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "users_email_key":
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}
