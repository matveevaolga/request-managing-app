package repository

import (
	"context"

	"github.com/matveevaolga/request-managing-app/internal/domain"
)

type ProjectTypeRepository interface {
	Create(ctx context.Context, projectType *domain.ProjectType) error
	GetByID(ctx context.Context, id int64) (*domain.ProjectType, error)
	GetByName(ctx context.Context, name string) (*domain.ProjectType, error)
	GetAllProjects(ctx context.Context) ([]domain.ProjectType, error)
	Update(ctx context.Context, projectType *domain.ProjectType) error
	Delete(ctx context.Context, id int64) error
}
