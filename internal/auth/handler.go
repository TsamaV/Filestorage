package auth

import (
	"context"
	"encoding/json"
	pkg "filestorage/pkg/response"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type AuthServiceInterface interface {
	SignUp(ctx context.Context, email, password string) (int, error)
	SignIn(ctx context.Context, email, password string) (string, error)
}

type Handler struct {
	AuthService AuthServiceInterface
	Logger      *zap.Logger
	validate    *validator.Validate
}

func NewHandler(authService AuthServiceInterface, logger *zap.Logger, validate *validator.Validate) *Handler {
	return &Handler{
		AuthService: authService,
		Logger:      logger,
		validate: validate,
	}
}

func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /sign-up", h.SignUp)
	mux.HandleFunc("POST /sign-in", h.SignIn)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.Logger.Info("signup request", zap.String("email", req.Email))

	userID, err := h.AuthService.SignUp(r.Context(), req.Email, req.Password)
	if err != nil {
		h.Logger.Error("signup error:", zap.Error(err))
		if err.Error() == "email already exists" {
			pkg.WriteError(w, http.StatusConflict, err.Error())
			return
		}
		pkg.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	pkg.WriteJSON(w, http.StatusCreated, SignUpResponse{
		ID:    userID,
		Email: req.Email,
	})

	h.Logger.Info("user created", zap.Int("user_id", userID))
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "invalid body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.Logger.Info("signin request", zap.String("email", req.Email))

	token, err := h.AuthService.SignIn(r.Context(), req.Email, req.Password)
	if err != nil {
		h.Logger.Error("sigin error:", zap.Error(err))
		pkg.WriteError(w, http.StatusUnauthorized, "internal error")
		return
	}

	pkg.WriteJSON(w, http.StatusOK, token)

	h.Logger.Info("user logged in", zap.String("email", req.Email))
}
