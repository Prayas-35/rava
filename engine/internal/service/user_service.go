package service

import (
	"context"
	"errors"
	"time"

	"github.com/Prayas-35/ragkit/engine/internal/database"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	exists := database.DB.QueryRow(ctx, `SELECT 1 FROM users WHERE email = $1`, req.Email)
	if err := exists.Scan(new(int)); err == nil {
		return nil, errors.New("email already in use")
	}

	row := database.DB.QueryRow(ctx,
		`INSERT INTO users (email, name, password_hash) VALUES ($1, $2, $3)
		RETURNING id, email, COALESCE(name, ''), created_at`,
		req.Email,
		req.Name,
		string(passwordHash),
	)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt); err != nil {
		return nil, err
	}

	return &u, nil
}

func AuthenticateUser(ctx context.Context, email string, password string) (*User, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT id, email, COALESCE(name, ''), created_at, password_hash
		 FROM users WHERE email = $1`,
		email,
	)

	var u User
	var passwordHash string
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.CreatedAt, &passwordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &u, nil
}
