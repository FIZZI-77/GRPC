package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"log"
	"os"
	"sso/internal/domain/models"
	"sso/internal/pkg/pgxhelper"
	"sso/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func StorageConnect(db *sql.DB) (*Storage, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	db, err := pgxhelper.NewPostgresDB(pgxhelper.Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSLMODE"),
	})
	if err != nil {
		log.Fatalf("Error connecting to postgres %v", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		log.Fatalf("Error closing postgres connection %v", err)
	}
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	const addUserQuery = `INSERT INTO users(email, pass_Hash) VALUES ($1, $2)`

	stmt, err := s.db.PrepareContext(ctx, addUserQuery)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.GetUserByEmail"

	const getUserByEmailQuery = `SELECT * FROM users WHERE email = $1`

	stmt, err := s.db.PrepareContext(ctx, getUserByEmailQuery)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	res := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = res.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}
