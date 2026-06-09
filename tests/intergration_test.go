package tests

import (
	"context"
	"encoding/json"
	"filestorage/internal/auth"
	"filestorage/internal/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func TestFullCycle(t *testing.T) {
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:password@localhost:5433/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	logger, _ := zap.NewDevelopment()
	validate := validator.New()

	userRepo := user.NewUserRepository(pool)
	authService := auth.NewAuthService(userRepo, "test-secret")
	authHandler := auth.NewHandler(authService, logger, validate)

	mux := http.NewServeMux()
	authHandler.Routes(mux)
	server := httptest.NewServer(mux)
	defer server.Close()

	// sign-up
	resp, _ := http.Post(server.URL+"/sign-up", "application/json", strings.NewReader(`{"email":"integration@test.com","password":"123456"}`))
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("sign-up: expected 201, got %d", resp.StatusCode)
	}

	// sign-in
	resp, _ = http.Post(server.URL+"/sign-in", "application/json", strings.NewReader(`{"email":"integration@test.com","password":"123456"}`))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("sign-up: expected 200, got %d", resp.StatusCode)
	}

	var token string
	json.NewDecoder(resp.Body).Decode(&token)

	if token == "" {
		t.Fatal("expected token, got empty")
	}

	t.Logf("toke: %s", token)
	t.Logf("full cycle passed")

	// upload
	
}
