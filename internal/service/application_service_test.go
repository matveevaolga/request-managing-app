package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository/mocks"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler/dto"
	"github.com/stretchr/testify/assert"
)

func TestApplicationService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)

	req := &dto.CreateApplicationRequest{
		FullName:         "Jane Smith",
		Email:            "jane@example.com",
		OrganisationName: "Startup Inc",
		ProjectName:      "Mobile App",
		TypeID:           1,
		ExpectedResults:  "iOS and Android MVP",
		IsPayed:          false,
	}

	t.Run("successful creation", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&domain.ProjectType{ID: 1, Name: "Test"}, nil)
		mockAppRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, app *domain.Application) error {
			app.ID = 123
			return nil
		})

		id, err := service.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, int64(123), id)
	})

	t.Run("project type not found", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(nil, domain.ErrProjectTypeNotFound)

		id, err := service.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, int64(0), id)
		assert.Contains(t, err.Error(), "invalid project type")
	})

	t.Run("duplicate application - returns error from repository", func(t *testing.T) {
		mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&domain.ProjectType{ID: 1, Name: "Test"}, nil)
		mockAppRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(domain.ErrApplicationAlreadyExists)

		id, err := service.Create(context.Background(), req)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrApplicationAlreadyExists, err)
		assert.Equal(t, int64(0), id)
	})
}

func TestApplicationService_Accept(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)

	t.Run("successful accept", func(t *testing.T) {
		app := &domain.Application{
			ID:     1,
			Status: domain.StatusPending,
		}
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)
		mockAppRepo.EXPECT().UpdateStatus(gomock.Any(), int64(1), domain.StatusAccepted, gomock.Any(), gomock.Any()).Return(nil)

		err := service.Accept(context.Background(), 1, 100)

		assert.NoError(t, err)
		assert.Equal(t, domain.StatusAccepted, app.Status)
	})

	t.Run("application not found", func(t *testing.T) {
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, domain.ErrApplicationNotFound)

		err := service.Accept(context.Background(), 999, 100)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrApplicationNotFound, err)
	})

	t.Run("application not pending", func(t *testing.T) {
		app := &domain.Application{
			ID:     1,
			Status: domain.StatusAccepted,
		}
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)

		err := service.Accept(context.Background(), 1, 100)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrApplicationNotPending, err)
	})
}

func TestApplicationService_Reject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)

	t.Run("successful reject", func(t *testing.T) {
		app := &domain.Application{
			ID:     1,
			Status: domain.StatusPending,
		}
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)
		mockAppRepo.EXPECT().UpdateStatus(gomock.Any(), int64(1), domain.StatusRejected, gomock.Any(), gomock.Any()).Return(nil)

		err := service.Reject(context.Background(), 1, 100, "test reason")

		assert.NoError(t, err)
		assert.Equal(t, domain.StatusRejected, app.Status)
	})

	t.Run("application not found", func(t *testing.T) {
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, domain.ErrApplicationNotFound)

		err := service.Reject(context.Background(), 999, 100, "reason")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrApplicationNotFound, err)
	})

	t.Run("application not pending", func(t *testing.T) {
		app := &domain.Application{
			ID:     1,
			Status: domain.StatusRejected,
		}
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(app, nil)

		err := service.Reject(context.Background(), 1, 100, "reason")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrApplicationNotPending, err)
	})
}

func TestApplicationService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	service := NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)

	t.Run("successful get", func(t *testing.T) {
		expected := &domain.Application{ID: 1, ProjectName: "Test"}
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(expected, nil)

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("application not found", func(t *testing.T) {
		mockAppRepo.EXPECT().GetByID(gomock.Any(), int64(999)).Return(nil, domain.ErrApplicationNotFound)

		result, err := service.GetByID(context.Background(), 999)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrApplicationNotFound, err)
		assert.Nil(t, result)
	})
}
