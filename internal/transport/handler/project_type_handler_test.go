package handler

import (
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
)

func TestProjectTypeHandler_GetAllProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	typeService := service.NewProjectTypeService(mockTypeRepo)
	handler := NewProjectTypeHandler(typeService)

	t.Run("successful get all project types", func(t *testing.T) {
		expected := []domain.ProjectType{
			{ID: 1, Name: "Startup"},
			{ID: 2, Name: "Internal"},
			{ID: 3, Name: "Personal"},
			{ID: 4, Name: "Social"},
		}
		mockTypeRepo.EXPECT().GetAllProjects(gomock.Any()).Return(expected, nil)

		req := httptest.NewRequest("GET", "/project/type", nil)
		rr := httptest.NewRecorder()

		handler.GetAllProjects(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp []dto.ProjectTypeResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 4)
		assert.Equal(t, "Startup", resp[0].Name)
	})

	t.Run("repository error", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetAllProjects(gomock.Any()).Return(nil, assert.AnError)

		req := httptest.NewRequest("GET", "/project/type", nil)
		rr := httptest.NewRecorder()

		handler.GetAllProjects(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var resp map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to get project types", resp["error"])
	})

	t.Run("empty list", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetAllProjects(gomock.Any()).Return([]domain.ProjectType{}, nil)

		req := httptest.NewRequest("GET", "/project/type", nil)
		rr := httptest.NewRecorder()

		handler.GetAllProjects(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp []dto.ProjectTypeResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Empty(t, resp)
	})
}
