package service

import (
	"context"
	"fmt"

	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository"
)

type ProjectTypeService struct {
	repo repository.ProjectTypeRepository
}

func NewProjectTypeService(repo repository.ProjectTypeRepository) *ProjectTypeService {
	return &ProjectTypeService{repo: repo}
}

func (s *ProjectTypeService) GetAll(ctx context.Context) ([]domain.ProjectType, error) {
	types, err := s.repo.GetAllProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get project types: %w", err)
	}
	return types, nil
}

func (s *ProjectTypeService) GetByID(ctx context.Context, id int64) (*domain.ProjectType, error) {
	pt, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return pt, nil
}
