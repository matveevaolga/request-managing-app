package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository/mocks"
	"github.com/matveevaolga/request-managing-app/internal/service"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler/dto"
	"github.com/stretchr/testify/assert"
)

func TestApplicationHandler_GetAllFiltered_SortingAndFiltering(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	handler := NewApplicationHandler(appService)

	tests := []struct {
		name           string
		queryParams    string
		expectedParams repository.ApplicationFilterParameters
		expectedStatus int
	}{
		{
			name:        "filter by active = true (only PENDING)",
			queryParams: "?active=true",
			expectedParams: repository.ApplicationFilterParameters{
				Active:            boolPtr(true),
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "filter by active = false (exclude PENDING)",
			queryParams: "?active=false",
			expectedParams: repository.ApplicationFilterParameters{
				Active:            boolPtr(false),
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "search by project name",
			queryParams: "?search=Mobile",
			expectedParams: repository.ApplicationFilterParameters{
				Search:            strPtr("Mobile"),
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "filter by project type ID",
			queryParams: "?projectTypeId=2",
			expectedParams: repository.ApplicationFilterParameters{
				ProjectTypeID:     int64Ptr(2),
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "sort by updated_at ASC",
			queryParams: "?sortByDateUpdated=ASC",
			expectedParams: repository.ApplicationFilterParameters{
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "ASC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "sort by updated_at DESC",
			queryParams: "?sortByDateUpdated=DESC",
			expectedParams: repository.ApplicationFilterParameters{
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "pagination with limit and offset",
			queryParams: "?limit=10&offset=5",
			expectedParams: repository.ApplicationFilterParameters{
				Limit:             10,
				Offset:            5,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "combined filters: active + search + type + sort + pagination",
			queryParams: "?active=true&search=CRM&projectTypeId=1&sortByDateUpdated=ASC&limit=25&offset=10",
			expectedParams: repository.ApplicationFilterParameters{
				Active:            boolPtr(true),
				Search:            strPtr("CRM"),
				ProjectTypeID:     int64Ptr(1),
				Limit:             25,
				Offset:            10,
				SortByDateUpdated: "ASC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "empty params - default values",
			queryParams: "",
			expectedParams: repository.ApplicationFilterParameters{
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid limit - uses default",
			queryParams: "?limit=-5",
			expectedParams: repository.ApplicationFilterParameters{
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid offset - uses default",
			queryParams: "?offset=-10",
			expectedParams: repository.ApplicationFilterParameters{
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid sort parameter - uses default DESC",
			queryParams: "?sortByDateUpdated=INVALID",
			expectedParams: repository.ApplicationFilterParameters{
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "invalid project type ID - ignored",
			queryParams: "?projectTypeId=invalid",
			expectedParams: repository.ApplicationFilterParameters{
				Limit:             20,
				Offset:            0,
				SortByDateUpdated: "DESC",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			previews := []domain.ApplicationPreview{
				{ID: 1, ProjectName: "Test", Status: domain.StatusPending},
			}

			mockAppRepo.EXPECT().
				GetAllFiltered(gomock.Any(), gomock.Any()).
				DoAndReturn(func(ctx interface{}, params repository.ApplicationFilterParameters) ([]domain.ApplicationPreview, int, error) {
					assert.Equal(t, tt.expectedParams.Limit, params.Limit)
					assert.Equal(t, tt.expectedParams.Offset, params.Offset)
					assert.Equal(t, tt.expectedParams.SortByDateUpdated, params.SortByDateUpdated)

					if tt.expectedParams.Active != nil {
						assert.NotNil(t, params.Active)
						assert.Equal(t, *tt.expectedParams.Active, *params.Active)
					} else {
						assert.Nil(t, params.Active)
					}

					if tt.expectedParams.Search != nil {
						assert.NotNil(t, params.Search)
						assert.Equal(t, *tt.expectedParams.Search, *params.Search)
					} else {
						assert.Nil(t, params.Search)
					}

					if tt.expectedParams.ProjectTypeID != nil {
						assert.NotNil(t, params.ProjectTypeID)
						assert.Equal(t, *tt.expectedParams.ProjectTypeID, *params.ProjectTypeID)
					} else {
						assert.Nil(t, params.ProjectTypeID)
					}

					return previews, len(previews), nil
				})

			req := httptest.NewRequest("GET", "/project/application/external/list"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			handler.GetAllFiltered(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestApplicationHandler_GetAllFiltered_ResponseStructure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	handler := NewApplicationHandler(appService)

	expectedPreviews := []domain.ApplicationPreview{
		{
			ID:               1,
			ProjectName:      "Mobile App",
			TypeName:         "Startup",
			Initiator:        "John Doe",
			OrganisationName: "Tech Corp",
			DateUpdated:      time.Now(),
			Status:           domain.StatusPending,
			RejectionMessage: nil,
		},
		{
			ID:               2,
			ProjectName:      "CRM System",
			TypeName:         "Internal",
			Initiator:        "Jane Smith",
			OrganisationName: "Business Inc",
			DateUpdated:      time.Now().Add(-24 * time.Hour),
			Status:           domain.StatusAccepted,
			RejectionMessage: nil,
		},
	}

	mockAppRepo.EXPECT().
		GetAllFiltered(gomock.Any(), gomock.Any()).
		Return(expectedPreviews, 2, nil)

	req := httptest.NewRequest("GET", "/project/application/external/list", nil)
	rr := httptest.NewRecorder()

	handler.GetAllFiltered(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp dto.ApplicationListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, 2, resp.Count)
	assert.Len(t, resp.Applications, 2)

	assert.Equal(t, int64(1), resp.Applications[0].ExternalApplicationID)
	assert.Equal(t, "Mobile App", resp.Applications[0].ProjectName)
	assert.Equal(t, "Startup", resp.Applications[0].TypeName)
	assert.Equal(t, "John Doe", resp.Applications[0].Initiator)
	assert.Equal(t, "Tech Corp", resp.Applications[0].OrganisationName)
	assert.Equal(t, string(domain.StatusPending), resp.Applications[0].Status)

	assert.Equal(t, int64(2), resp.Applications[1].ExternalApplicationID)
	assert.Equal(t, "CRM System", resp.Applications[1].ProjectName)
	assert.Equal(t, string(domain.StatusAccepted), resp.Applications[1].Status)
}

func TestApplicationHandler_GetAllFiltered_EmptyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	handler := NewApplicationHandler(appService)

	mockAppRepo.EXPECT().
		GetAllFiltered(gomock.Any(), gomock.Any()).
		Return([]domain.ApplicationPreview{}, 0, nil)

	req := httptest.NewRequest("GET", "/project/application/external/list?active=true", nil)
	rr := httptest.NewRecorder()

	handler.GetAllFiltered(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp dto.ApplicationListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Count)
	assert.Empty(t, resp.Applications)
}

func boolPtr(b bool) *bool {
	return &b
}

func int64Ptr(i int64) *int64 {
	return &i
}
