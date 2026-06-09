package auth

import (
	"context"
	"errors"
	"filestorage/internal/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	CreateUser(ctx context.Context, email string, passwordHash string) (int, error)
}

type AuthService struct {
	userRepo  UserRepo
	jwtSecret string
}

func NewAuthService(repo UserRepo, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

func (service *AuthService) SignUp(ctx context.Context, email, password string) (int, error) {
	_, err := service.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return 0, errors.New("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	userID, err := service.userRepo.CreateUser(ctx, email, string(hash))
	if err != nil {
		return 0, err
	}

	return userID, err
}

func (service *AuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := service.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(service.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, err
}
