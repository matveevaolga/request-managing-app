package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/service"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler/dto"
)

type AuthHandler struct {
	authService *service.AuthService
	validate    *validator.Validate
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validator.New(),
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error(), err)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			RespondWithError(w, http.StatusUnauthorized, "Invalid credentials", err)
		default:
			RespondWithError(w, http.StatusInternalServerError, "Failed to login", err)
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, dto.LoginResponse{Token: token})
}
