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

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func StorageConnect() (*sql.DB, error) {
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

	return db, nil
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		log.Fatalf("Error closing postgres connection %v", err)
	}
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	const addUserQuery = `INSERT INTO users(email, pass_Hash) VALUES ($1, $2) RETURNING id`

	var id int64
	stmt, err := s.db.PrepareContext(ctx, addUserQuery)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = stmt.QueryRowContext(ctx, email, passHash).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

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

	var isAdmin bool
	
	var user models.User
	err = res.Scan(&user.ID, &user.Email, &user.PassHash, &isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	isAdminQuery := `SELECT is_admin FROM users WHERE id = $1`
	stmt, err := s.db.PrepareContext(ctx, isAdminQuery)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)
	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}

func (s *Storage) GetAppByID(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgres.GetAppByID"
	const getAppByIDQuery = `SELECT * FROM apps WHERE id = $1`

	stmt, err := s.db.PrepareContext(ctx, getAppByIDQuery)
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	var app models.App

	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
