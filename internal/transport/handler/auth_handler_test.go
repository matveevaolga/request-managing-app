package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository/mocks"
	"github.com/matveevaolga/request-managing-app/internal/service"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler/dto"
	"github.com/matveevaolga/request-managing-app/internal/transport/middleware"
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

func TestAuthorization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)

	authService := service.NewAuthService(mockUserRepo, "test-secret", 24)
	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	typeService := service.NewProjectTypeService(mockTypeRepo)

	authHandler := NewAuthHandler(authService)
	appHandler := NewApplicationHandler(appService)
	typeHandler := NewProjectTypeHandler(typeService)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	adminUser := &domain.User{
		ID:       1,
		Username: "admin",
		Password: string(hashedPassword),
		Role:     domain.RoleAdmin,
	}

	regularUser := &domain.User{
		ID:       2,
		Username: "user",
		Password: string(hashedPassword),
		Role:     domain.RoleUser,
	}

	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		setupMocks     func()
		user           *domain.User
		expectedStatus int
	}{
		{
			name:   "POST /login - public",
			method: "POST",
			path:   "/login",
			body:   dto.LoginRequest{Login: "admin", Password: "password"},
			setupMocks: func() {
				mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin").Return(adminUser, nil)
			},
			user:           nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET /project/type - public",
			method: "GET",
			path:   "/project/type",
			body:   nil,
			setupMocks: func() {
				mockTypeRepo.EXPECT().GetAllProjects(gomock.Any()).Return([]domain.ProjectType{
					{ID: 1, Name: "Startup"},
				}, nil)
			},
			user:           nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "POST /project/application/external - public",
			method: "POST",
			path:   "/project/application/external",
			body:   dto.CreateApplicationRequest{FullName: "Test", Email: "test@test.com", OrganisationName: "Org", ProjectName: "Proj", TypeID: 1, ExpectedResults: "Res"},
			setupMocks: func() {
				mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&domain.ProjectType{ID: 1, Name: "Test"}, nil)
				mockAppRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, app *domain.Application) error {
					app.ID = 123
					return nil
				})
			},
			user:           nil,
			expectedStatus: http.StatusOK,
		},

		{
			name:   "GET /project/application/external/list - admin success",
			method: "GET",
			path:   "/project/application/external/list",
			body:   nil,
			setupMocks: func() {
				mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin").Return(adminUser, nil)
				mockAppRepo.EXPECT().GetAllFiltered(gomock.Any(), gomock.Any()).Return([]domain.ApplicationPreview{}, 0, nil)
			},
			user:           adminUser,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET /project/application/external/list - user forbidden",
			method: "GET",
			path:   "/project/application/external/list",
			body:   nil,
			setupMocks: func() {
				mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "user").Return(regularUser, nil)
			},
			user:           regularUser,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "GET /project/application/external/list - no auth",
			method:         "GET",
			path:           "/project/application/external/list",
			body:           nil,
			setupMocks:     func() {},
			user:           nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "POST /project/application/external/1/accept - admin success",
			method: "POST",
			path:   "/project/application/external/1/accept",
			body:   nil,
			setupMocks: func() {
				app := &domain.Application{ID: 1, Status: domain.StatusPending}
				mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin").Return(adminUser, nil)
				mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)
				mockAppRepo.EXPECT().UpdateStatus(gomock.Any(), int64(1), domain.StatusAccepted, gomock.Any(), gomock.Any()).Return(nil)
			},
			user:           adminUser,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "POST /project/application/external/1/accept - user forbidden",
			method: "POST",
			path:   "/project/application/external/1/accept",
			body:   nil,
			setupMocks: func() {
				mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "user").Return(regularUser, nil)
			},
			user:           regularUser,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "POST /project/application/external/1/accept - no auth",
			method:         "POST",
			path:           "/project/application/external/1/accept",
			body:           nil,
			setupMocks:     func() {},
			user:           nil,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "POST /project/application/external/1/reject - admin success",
			method: "POST",
			path:   "/project/application/external/1/reject",
			body:   dto.RejectRequest{Reason: "test"},
			setupMocks: func() {
				app := &domain.Application{ID: 1, Status: domain.StatusPending}
				mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin").Return(adminUser, nil)
				mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)
				mockAppRepo.EXPECT().UpdateStatus(gomock.Any(), int64(1), domain.StatusRejected, gomock.Any(), gomock.Any()).Return(nil)
			},
			user:           adminUser,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET /project/application/external/1 - admin success",
			method: "GET",
			path:   "/project/application/external/1",
			body:   nil,
			setupMocks: func() {
				app := &domain.Application{ID: 1, FullName: "Test", Email: "test@test.com", OrganisationName: "Org", ProjectName: "Proj", TypeID: 1, ExpectedResults: "Res", Status: domain.StatusPending}
				mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin").Return(adminUser, nil)
				mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)
				mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&domain.ProjectType{ID: 1, Name: "Test"}, nil)
			},
			user:           adminUser,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET /project/application/external/1 - user forbidden",
			method: "GET",
			path:   "/project/application/external/1",
			body:   nil,
			setupMocks: func() {
				mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "user").Return(regularUser, nil)
			},
			user:           regularUser,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "GET /project/application/external/1 - no auth",
			method:         "GET",
			path:           "/project/application/external/1",
			body:           nil,
			setupMocks:     func() {},
			user:           nil,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			var req *http.Request
			if tt.body != nil {
				jsonBody, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewReader(jsonBody))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}

			if tt.user != nil {
				token, err := authService.Login(context.Background(), tt.user.Username, "password")
				assert.NoError(t, err)
				req.Header.Set("X-API-TOKEN", token)
			}

			mux := setupRouter(authHandler, typeHandler, appHandler, authService)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func setupRouter(authHandler *AuthHandler, typeHandler *ProjectTypeHandler, appHandler *ApplicationHandler, authService *service.AuthService) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("GET /project/type", typeHandler.GetAllProjects)

	mux.HandleFunc("POST /project/application/external", appHandler.Create)

	mux.HandleFunc("GET /project/application/external/list", middleware.CheckAdmin(appHandler.GetAllFiltered, authService))
	mux.HandleFunc("GET /project/application/external/{applicationId}", middleware.CheckAdmin(appHandler.GetByID, authService))
	mux.HandleFunc("POST /project/application/external/{applicationId}/accept", middleware.CheckAdmin(appHandler.Accept, authService))
	mux.HandleFunc("POST /project/application/external/{applicationId}/reject", middleware.CheckAdmin(appHandler.Reject, authService))

	return mux
}
