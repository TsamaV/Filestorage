package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type MockAuthService struct{}

func (m *MockAuthService) SignUp(ctx context.Context, email, password string) (int, error) {
	return 1, nil
}

func (m *MockAuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	return "token", nil
}

func TestSignUp(t *testing.T) {
	body := `{"email":"test@test.com", "password":"123456"}`

	req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(body))
	w := httptest.NewRecorder()

	mock := &MockAuthService{}
	logger, _ := zap.NewDevelopment()
	validate := validator.New()

	handler := NewHandler(mock, logger, validate)

	handler.SignUp(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestSignIn(t *testing.T) {
	body := `{"email":"test@test.com", "password":"123456"}`

	req := httptest.NewRequest(http.MethodPost, "/sign-in", strings.NewReader(body))
	w := httptest.NewRecorder()

	mock := &MockAuthService{}
	logger, _ := zap.NewDevelopment()
	validate := validator.New()

	handler := NewHandler(mock, logger, validate)
	handler.SignIn(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestBadSignUp(t *testing.T) {
	body := `{"email":"", "password":"1234"}`

	req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(body))
	w := httptest.NewRecorder()

	mock := &MockAuthService{}
	logger, _ := zap.NewDevelopment()
	validate := validator.New()

	handler := NewHandler(mock, logger, validate)

	handler.SignUp(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestBadSignIn(t *testing.T) {
	body := `{"email":"test@test.com", "password":""}`

	req := httptest.NewRequest(http.MethodPost, "/sign-in", strings.NewReader(body))
	w := httptest.NewRecorder()

	mock := &MockAuthService{}
	logger, _ := zap.NewDevelopment()
	validate := validator.New()

	handler := NewHandler(mock, logger, validate)
	handler.SignIn(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusOK, w.Code)
	}
}