package auth

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	GetAppByID(ctx context.Context, appID int) (models.App, error)
}

func NewAuth(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, token time.Duration) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     token,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "auth.Login"
	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("attempting to login user")

	user, err := a.userProvider.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		a.log.Error("failed to get user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.GetAppByID(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("successfully logged in")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to create token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (int64, error) {
	const op = "auth.RegisterNewUser"
	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", err)
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user created")

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(slog.String("op", op), slog.Int64("user_id", userID))

	log.Info("checking if user exists")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", err)

			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("checking if user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}
