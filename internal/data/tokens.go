package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zrotrasukha/jobman/internal/validator"
)

var (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

type TokenModel struct {
	pool *pgxpool.Pool
}

type TokenModelInterface interface {
	New(userID int64, ttl time.Duration, scope string) (*Token, error)
	Insert(token *Token) error
	DeleteAllforUser(userID int64, scope string) error
}

func ValidateTokenPlainText(v *validator.Validator, tokenPlainText string) {
	v.CheckField(tokenPlainText != "", "token", "must be provided")
	v.CheckField(len(tokenPlainText) == 26, "token", "must be 26 bytes long")
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, users_id, expiry, scope) VALUES ($1, $2, $3, $4)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.pool.Exec(ctx, query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (m TokenModel) DeleteAllforUser(userID int64, scope string) error {
	query := `DELETE FROM tokens WHERE users_id = $1 AND scope = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.pool.Exec(ctx, query, userID, scope)
	return err
}

func GenerateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	//NOTE TO SELF: randomBytes has total of 16 * 8 = 128 bits, and base32 encoding reads 5 bits for each character,
	// so the resulting string will be 128 / 5 = 25.6 characters long, which is rounded up to 26 characters.
	// We don't want any padding because it makes the token human unfriendly.
	encodedString := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	token.Plaintext = encodedString

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}
