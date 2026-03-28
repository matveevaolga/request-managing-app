package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestProjectTypeService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	service := NewProjectTypeService(mockTypeRepo)

	t.Run("successful get all", func(t *testing.T) {
		expected := []domain.ProjectType{
			{ID: 1, Name: "Startup"},
			{ID: 2, Name: "Internal"},
			{ID: 3, Name: "Personal"},
			{ID: 4, Name: "Social"},
		}
		mockTypeRepo.EXPECT().GetAllProjects(gomock.Any()).Return(expected, nil)

		result, err := service.GetAllProjects(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		assert.Len(t, result, 4)
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetAllProjects(gomock.Any()).Return(nil, domain.ErrProjectTypeNotFound)

		result, err := service.GetAllProjects(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get project types")
	})

	t.Run("empty list", func(t *testing.T) {
		expected := []domain.ProjectType{}
		mockTypeRepo.EXPECT().GetAllProjects(gomock.Any()).Return(expected, nil)

		result, err := service.GetAllProjects(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestProjectTypeService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	service := NewProjectTypeService(mockTypeRepo)

	t.Run("successful get by id", func(t *testing.T) {
		expected := &domain.ProjectType{ID: 1, Name: "Startup"}
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(expected, nil)

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("project type not found", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, domain.ErrProjectTypeNotFound)

		result, err := service.GetByID(context.Background(), 999)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrProjectTypeNotFound, err)
		assert.Nil(t, result)
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(nil, assert.AnError)

		result, err := service.GetByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
