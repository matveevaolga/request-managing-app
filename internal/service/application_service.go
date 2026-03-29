package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository"
	"github.com/matveevaolga/request-managing-app/internal/transport/handler/dto"
)

type ApplicationService struct {
	appRepo  repository.ApplicationRepository
	typeRepo repository.ProjectTypeRepository
	userRepo repository.UserRepository
}

func NewApplicationService(appRepo repository.ApplicationRepository,
	typeRepo repository.ProjectTypeRepository, userRepo repository.UserRepository,
) *ApplicationService {
	return &ApplicationService{
		appRepo:  appRepo,
		typeRepo: typeRepo,
		userRepo: userRepo,
	}
}

func (s *ApplicationService) Create(ctx context.Context, req *dto.CreateApplicationRequest) (int64, error) {
	if req.Phone != nil && *req.Phone != "" {
		pattern := `.*\+7\s\(\d{3}\)\s\d{3}-\d{2}-\d{2}.*`
		if matched, _ := regexp.MatchString(pattern, *req.Phone); !matched {
			return 0, fmt.Errorf("invalid phone format. Expected format: +7 (XXX) XXX-XX-XX")
		}
	}
	_, err := s.typeRepo.GetByID(ctx, req.TypeID)
	if err != nil {
		return 0, fmt.Errorf("invalid project type: %w", err)
	}
	app := &domain.Application{
		FullName:              req.FullName,
		Email:                 req.Email,
		Phone:                 req.Phone,
		OrganisationName:      req.OrganisationName,
		OrganisationURL:       req.OrganisationURL,
		ProjectName:           req.ProjectName,
		TypeID:                req.TypeID,
		ExpectedResults:       req.ExpectedResults,
		IsPayed:               req.IsPayed,
		AdditionalInformation: req.AdditionalInformation,
		Status:                domain.StatusPending,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
	err = s.appRepo.Create(ctx, app)
	if err != nil {
		if errors.Is(err, domain.ErrApplicationAlreadyExists) {
			return 0, domain.ErrApplicationAlreadyExists
		}
		return 0, fmt.Errorf("failed to create application: %w", err)
	}
	return app.ID, nil
}

func (s *ApplicationService) GetByID(ctx context.Context, id int64) (*domain.Application, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (s *ApplicationService) GetAllFiltered(ctx context.Context, params repository.ApplicationFilterParameters) ([]domain.ApplicationPreview, int, error) {
	apps, total, err := s.appRepo.GetAllFiltered(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get applications: %w", err)
	}
	return apps, total, nil
}

func (s *ApplicationService) Accept(ctx context.Context, id int64, reviewerID int64) error {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !app.IsPending() {
		return domain.ErrApplicationNotPending
	}

	err = app.Accept(reviewerID)
	if err != nil {
		return err
	}

	return s.appRepo.UpdateStatus(ctx, id, app.Status, app.ReviewerID, app.RejectedReason)
}

func (s *ApplicationService) Reject(ctx context.Context, id int64, reviewerID int64, reason string) error {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !app.IsPending() {
		return domain.ErrApplicationNotPending
	}

	err = app.Reject(reviewerID, reason)
	if err != nil {
		return err
	}

	return s.appRepo.UpdateStatus(ctx, id, app.Status, app.ReviewerID, app.RejectedReason)
}

func (s *ApplicationService) GetProjectTypeByID(ctx context.Context, id int64) (*domain.ProjectType, error) {
	return s.typeRepo.GetByID(ctx, id)
}
