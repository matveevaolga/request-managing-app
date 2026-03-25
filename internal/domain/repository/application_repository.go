package repository

import (
	"context"

	"github.com/matveevaolga/request-managing-app/internal/domain"
)

type ApplicationFilterParameters struct {
	Active            *bool
	Search            *string
	ProjectTypeID     *int64
	SortByDateUpdated string
	Limit             int
	Offset            int
}

type ApplicationRepository interface {
	Create(ctx context.Context, app *domain.Application) error
	GetByID(ctx context.Context, id int64) (*domain.Application, error)
	GetAllFiltered(ctx context.Context, params ApplicationFilterParameters) ([]domain.ApplicationPreview, int, error)
	Update(ctx context.Context, app *domain.Application) error
	UpdateStatus(ctx context.Context, id int64, status domain.ApplicationStatus, reviewerID *int64, rejectedReason *string) error
}
