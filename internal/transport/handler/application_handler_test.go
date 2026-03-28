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
	"github.com/matveevaolga/request-managing-app/internal/domain/repository"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository/mocks"
	"github.com/matveevaolga/request-managing-app/internal/service"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler/dto"
	"github.com/matveevaolga/request-managing-app/internal/transport/middleware"
	"github.com/stretchr/testify/assert"
)

func TestApplicationHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	handler := NewApplicationHandler(appService)

	reqBody := dto.CreateApplicationRequest{
		FullName:         "Jane Smith",
		Email:            "jane@example.com",
		OrganisationName: "Startup Inc",
		ProjectName:      "Mobile App",
		TypeID:           1,
		ExpectedResults:  "iOS and Android MVP",
		IsPayed:          false,
	}

	t.Run("successful creation", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&domain.ProjectType{ID: 1, Name: "Startup"}, nil)
		mockAppRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx interface{}, app *domain.Application) error {
			app.ID = 123
			return nil
		})

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/project/application/external", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var id int64
		err := json.Unmarshal(rr.Body.Bytes(), &id)
		assert.NoError(t, err)
		assert.Equal(t, int64(123), id)
	})

	t.Run("invalid request body", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/project/application/external", bytes.NewReader([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("validation failed - missing required field", func(t *testing.T) {
		invalidReq := dto.CreateApplicationRequest{
			FullName: "Jane Smith",
		}
		jsonBody, _ := json.Marshal(invalidReq)
		req := httptest.NewRequest("POST", "/project/application/external", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("project type not found", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(nil, domain.ErrProjectTypeNotFound)

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/project/application/external", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var resp map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid project type", resp["error"])
	})

	t.Run("duplicate application", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&domain.ProjectType{ID: 1, Name: "Startup"}, nil)
		mockAppRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.ErrApplicationAlreadyExists)

		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/project/application/external", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		handler.Create(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var resp map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Application with this project name and email already exists", resp["error"])
	})
}

func TestApplicationHandler_GetAllFiltered(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	handler := NewApplicationHandler(appService)

	t.Run("successful get list", func(t *testing.T) {
		previews := []domain.ApplicationPreview{
			{ID: 1, ProjectName: "Test Project", TypeName: "Startup", Initiator: "Test User", OrganisationName: "Test Org", Status: domain.StatusPending},
			{ID: 2, ProjectName: "Another Project", TypeName: "Internal", Initiator: "Another User", OrganisationName: "Another Org", Status: domain.StatusAccepted},
		}
		mockAppRepo.EXPECT().GetAllFiltered(gomock.Any(), gomock.Any()).Return(previews, 2, nil)

		req := httptest.NewRequest("GET", "/project/application/external/list?limit=20&offset=0", nil)
		rr := httptest.NewRecorder()

		handler.GetAllFiltered(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp dto.ApplicationListResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, 2, resp.Count)
		assert.Len(t, resp.Applications, 2)
	})

	t.Run("repository error", func(t *testing.T) {
		mockAppRepo.EXPECT().GetAllFiltered(gomock.Any(), gomock.Any()).Return(nil, 0, assert.AnError)

		req := httptest.NewRequest("GET", "/project/application/external/list", nil)
		rr := httptest.NewRecorder()

		handler.GetAllFiltered(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var resp map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to get applications", resp["error"])
	})

	t.Run("with filters", func(t *testing.T) {
		mockAppRepo.EXPECT().GetAllFiltered(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx interface{}, params repository.ApplicationFilterParameters) ([]domain.ApplicationPreview, int, error) {
			assert.Equal(t, true, *params.Active)
			assert.Equal(t, "Test", *params.Search)
			assert.Equal(t, int64(1), *params.ProjectTypeID)
			assert.Equal(t, "ASC", params.SortByDateUpdated)
			assert.Equal(t, 10, params.Limit)
			assert.Equal(t, 5, params.Offset)
			return []domain.ApplicationPreview{}, 0, nil
		})

		req := httptest.NewRequest("GET", "/project/application/external/list?active=true&search=Test&projectTypeId=1&sortByDateUpdated=ASC&limit=10&offset=5", nil)
		rr := httptest.NewRecorder()

		handler.GetAllFiltered(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestApplicationHandler_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	handler := NewApplicationHandler(appService)

	t.Run("successful get by id", func(t *testing.T) {
		app := &domain.Application{
			ID:               1,
			FullName:         "Jane Smith",
			Email:            "jane@example.com",
			OrganisationName: "Startup Inc",
			ProjectName:      "Mobile App",
			TypeID:           1,
			ExpectedResults:  "iOS and Android MVP",
			IsPayed:          false,
			Status:           domain.StatusPending,
		}
		projectType := &domain.ProjectType{ID: 1, Name: "Startup"}

		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(projectType, nil)

		req := httptest.NewRequest("GET", "/project/application/external/1", nil)
		rr := httptest.NewRecorder()

		handler.GetByID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp dto.ApplicationResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.ApplicationID)
		assert.Equal(t, "Jane Smith", resp.FullName)
		assert.Equal(t, "Startup", resp.TypeName)
	})

	t.Run("application not found", func(t *testing.T) {
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, domain.ErrApplicationNotFound)

		req := httptest.NewRequest("GET", "/project/application/external/999", nil)
		rr := httptest.NewRecorder()

		handler.GetByID(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)

		var resp map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Application not found", resp["error"])
	})
}

func TestApplicationHandler_Accept(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	handler := NewApplicationHandler(appService)

	t.Run("successful accept", func(t *testing.T) {
		app := &domain.Application{
			ID:     1,
			Status: domain.StatusPending,
		}
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)
		mockAppRepo.EXPECT().UpdateStatus(gomock.Any(), int64(1), domain.StatusAccepted, gomock.Any(), gomock.Any()).Return(nil)

		req := httptest.NewRequest("POST", "/project/application/external/1/accept", nil)

		ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(100))
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler.Accept(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("application not found", func(t *testing.T) {
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, domain.ErrApplicationNotFound)

		req := httptest.NewRequest("POST", "/project/application/external/999/accept", nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(100))
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.Accept(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("application not pending", func(t *testing.T) {
		app := &domain.Application{
			ID:     1,
			Status: domain.StatusAccepted,
		}
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)

		req := httptest.NewRequest("POST", "/project/application/external/1/accept", nil)
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(100))
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.Accept(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestApplicationHandler_Reject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	handler := NewApplicationHandler(appService)

	t.Run("successful reject", func(t *testing.T) {
		app := &domain.Application{
			ID:     1,
			Status: domain.StatusPending,
		}
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)
		mockAppRepo.EXPECT().UpdateStatus(gomock.Any(), int64(1), domain.StatusRejected, gomock.Any(), gomock.Any()).Return(nil)

		reqBody := dto.RejectRequest{Reason: "Test reason"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/project/application/external/1/reject", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(100))
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.Reject(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("missing reason", func(t *testing.T) {
		reqBody := dto.RejectRequest{}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/project/application/external/1/reject", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(100))
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.Reject(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("application not found", func(t *testing.T) {
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, domain.ErrApplicationNotFound)

		reqBody := dto.RejectRequest{Reason: "Test reason"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/project/application/external/999/reject", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(100))
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.Reject(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("application not pending", func(t *testing.T) {
		app := &domain.Application{
			ID:     1,
			Status: domain.StatusRejected,
		}
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)

		reqBody := dto.RejectRequest{Reason: "Test reason"}
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/project/application/external/1/reject", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), middleware.UserIDKey, int64(100))
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler.Reject(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
