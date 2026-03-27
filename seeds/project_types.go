package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/matveevaolga/request-managing-app/internal/domain"
	domainrepo "github.com/matveevaolga/request-managing-app/internal/domain/repository"
)

func seedProjectTypes(ctx context.Context, repo domainrepo.ProjectTypeRepository) error {
	types, err := loadProjectTypes()
	if err != nil {
		return err
	}
	for _, t := range types {
		if err := createProjectType(ctx, repo, t); err != nil {
			return err
		}
	}
	return nil
}

func loadProjectTypes() ([]typeSeed, error) {
	data, err := seedFS.ReadFile("data/project_types.json")
	if err != nil {
		return nil, err
	}
	var types []typeSeed
	if err := json.Unmarshal(data, &types); err != nil {
		return nil, err
	}
	return types, nil
}

func createProjectType(ctx context.Context, repo domainrepo.ProjectTypeRepository, t typeSeed) error {
	exists, err := repo.GetByName(ctx, t.Name)
	if err != nil && err != domain.ErrProjectTypeNotFound {
		return err
	}
	if exists != nil {
		slog.Info("Project type already exists", "name", t.Name)
		return nil
	}
	pt := &domain.ProjectType{
		Name:      t.Name,
		CreatedAt: time.Now(),
	}
	if err := repo.Create(ctx, pt); err != nil {
		return err
	}
	slog.Info("Project type created", "name", pt.Name, "id", pt.ID)
	return nil
}
