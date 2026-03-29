package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository/mocks"
	"github.com/matveevaolga/request-managing-app/internal/service"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler/dto"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	authService := service.NewAuthService(mockUserRepo, "test-secret", 24)
	handler := NewAuthHandler(authService)

	t.Run("successful login", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       1,
			Username: "admin1",
			Password: string(hashedPassword),
			Role:     domain.RoleAdmin,
		}
		mockUserRepo.EXPECT().
			GetByUsername(gomock.Any(), "admin1").
			Return(user, nil)

		reqBody := dto.LoginRequest{
			Login:    "admin1",
			Password: "admin1",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp dto.LoginResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("invalid credentials - user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().
			GetByUsername(gomock.Any(), "nonexistent").
			Return(nil, domain.ErrUserNotFound)

		reqBody := dto.LoginRequest{
			Login:    "nonexistent",
			Password: "password",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		var resp map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid credentials", resp["error"])
	})

	t.Run("wrong password", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:       1,
			Username: "admin1",
			Password: string(hashedPassword),
			Role:     domain.RoleAdmin,
		}
		mockUserRepo.EXPECT().
			GetByUsername(gomock.Any(), "admin1").
			Return(user, nil)

		reqBody := dto.LoginRequest{
			Login:    "admin1",
			Password: "wrongpassword",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		var resp map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid credentials", resp["error"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Login(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
