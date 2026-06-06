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
	Version   int32     `json:"-"`
}

type Password struct {
	PlainText *string
	Hash      []byte
}

func (p *Password) Set(PlainTextPassword string) error {
	b, err := bcrypt.GenerateFromPassword([]byte(PlainTextPassword), 12)
	if err != nil {
		return err
	}

	p.PlainText = &PlainTextPassword
	p.Hash = b

	return nil
}

func (p *Password) Matches(PlainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(PlainTextPassword))

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

	if user.Password.Hash == nil {
		panic("missing Password Hash for user")
	}
}

type UserModelInterface interface {
	Insert(user *User) error
}

type UserModel struct {
	pool *pgxpool.Pool
}

func (m UserModel) Insert(user *User) error {
	query := `INSERT INTO users (name, email, password_hash, activated)
						VALUES ($1, $2, $3, $4)
						RETURNING id, created_at, version`

	args := []any{user.Name, user.Email, user.Password.Hash, user.Activated}

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

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, created_at, name, email, Password_Hash, activated, version
						FROM users
						WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.pool.QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `UPDATE users
						SET name = $1, email = $2, Password_Hash = $3, activated = $4, version = version + 1
						WHERE id = $5 AND version = $6
						RETURNING version`

	args := []any{user.Name, user.Email, user.Password.Hash, user.Activated, user.Id, user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.pool.QueryRow(ctx, query, args...).Scan(&user.Version)
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

func (m UserModel) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
