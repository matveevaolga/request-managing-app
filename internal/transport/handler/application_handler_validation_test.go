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
)

func TestApplicationHandler_Validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAppRepo := mocks.NewMockApplicationRepository(ctrl)
	mockTypeRepo := mocks.NewMockProjectTypeRepository(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	appService := service.NewApplicationService(mockAppRepo, mockTypeRepo, mockUserRepo)
	handler := NewApplicationHandler(appService)

	tests := []struct {
		name       string
		request    dto.CreateApplicationRequest
		wantStatus int
	}{
		{
			name: "valid application",
			request: dto.CreateApplicationRequest{
				FullName:         "Some User",
				Email:            "email@example.com",
				Phone:            strPtr("+7 (999) 123-45-67"),
				OrganisationName: "Some Organisation",
				OrganisationURL:  strPtr("https://example.com"),
				ProjectName:      "Project Name",
				TypeID:           1,
				ExpectedResults:  "Some results",
				IsPayed:          false,
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid phone format",
			request: dto.CreateApplicationRequest{
				FullName:         "Some User",
				Email:            "email@example.com",
				Phone:            strPtr("invalid-phone"),
				OrganisationName: "Some Organisation",
				ProjectName:      "Project Name",
				TypeID:           1,
				ExpectedResults:  "Some results",
				IsPayed:          false,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email format",
			request: dto.CreateApplicationRequest{
				FullName:         "Some User",
				Email:            "not-an-email",
				OrganisationName: "Some Organisation",
				ProjectName:      "Project Name",
				TypeID:           1,
				ExpectedResults:  "Some results",
				IsPayed:          false,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid URL format",
			request: dto.CreateApplicationRequest{
				FullName:         "Some User",
				Email:            "email@example.com",
				OrganisationName: "Some Organisation",
				OrganisationURL:  strPtr("not-a-url"),
				ProjectName:      "Project Name",
				TypeID:           1,
				ExpectedResults:  "Some results",
				IsPayed:          false,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing full name",
			request: dto.CreateApplicationRequest{
				FullName:         "",
				Email:            "email@example.com",
				OrganisationName: "Some Organisation",
				ProjectName:      "Project Name",
				TypeID:           1,
				ExpectedResults:  "Some results",
				IsPayed:          false,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing organisation name",
			request: dto.CreateApplicationRequest{
				FullName:         "Some User",
				Email:            "email@example.com",
				OrganisationName: "",
				ProjectName:      "Project Name",
				TypeID:           1,
				ExpectedResults:  "Some results",
				IsPayed:          false,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing project name",
			request: dto.CreateApplicationRequest{
				FullName:         "Some User",
				Email:            "email@example.com",
				OrganisationName: "Some Organisation",
				ProjectName:      "",
				TypeID:           1,
				ExpectedResults:  "Some results",
				IsPayed:          false,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing expected results",
			request: dto.CreateApplicationRequest{
				FullName:         "Some User",
				Email:            "email@example.com",
				OrganisationName: "Some Organisation",
				ProjectName:      "Project Name",
				TypeID:           1,
				ExpectedResults:  "",
				IsPayed:          false,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantStatus == http.StatusOK {
				mockTypeRepo.EXPECT().GetByID(gomock.Any(), int64(1)).Return(&domain.ProjectType{ID: 1, Name: "Startup"}, nil)
				mockAppRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx interface{}, app *domain.Application) error {
					app.ID = 123
					return nil
				})
			}

			jsonBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/project/application/external", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.Create(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func strPtr(s string) *string {
	return &s
}
