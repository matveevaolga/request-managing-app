package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/matveevaolga/request-managing-app/internal/domain"
	domainrepo "github.com/matveevaolga/request-managing-app/internal/domain/repository"
)

func seedApplications(ctx context.Context, appRepo domainrepo.ApplicationRepository, typeRepo domainrepo.ProjectTypeRepository) error {
	types, err := typeRepo.GetAllProjects(ctx)
	if err != nil {
		return err
	}
	typeMap := buildTypeMap(types)
	apps, err := loadApplications()
	if err != nil {
		return err
	}
	for _, a := range apps {
		if err := createApplication(ctx, appRepo, typeMap, a); err != nil {
			return err
		}
	}
	return nil
}

func loadApplications() ([]applicationSeed, error) {
	data, err := seedFS.ReadFile("data/applications.json")
	if err != nil {
		return nil, err
	}
	var apps []applicationSeed
	if err := json.Unmarshal(data, &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

func buildTypeMap(types []domain.ProjectType) map[string]int64 {
	m := make(map[string]int64)
	for _, t := range types {
		m[t.Name] = t.ID
	}
	return m
}

func createApplication(ctx context.Context, repo domainrepo.ApplicationRepository, typeMap map[string]int64, a applicationSeed) error {
	typeID, ok := typeMap[a.TypeName]
	if !ok {
		slog.Warn("Project type not found", "type", a.TypeName)
		return nil
	}
	app := fillApplication(a, typeID)
	setApplicationDates(app, a)
	if err := repo.Create(ctx, app); err != nil {
		if errors.Is(err, domain.ErrApplicationAlreadyExists) {
			slog.Info("Application already exists", "project", app.ProjectName, "email", app.Email)
			return nil
		}
		return err
	}
	slog.Info("Application created", "id", app.ID, "project", app.ProjectName, "status", app.Status)
	return nil
}

func fillApplication(a applicationSeed, typeID int64) *domain.Application {
	return &domain.Application{
		FullName:              a.FullName,
		Email:                 a.Email,
		Phone:                 a.Phone,
		OrganisationName:      a.OrganisationName,
		OrganisationURL:       a.OrganisationURL,
		ProjectName:           a.ProjectName,
		TypeID:                typeID,
		ExpectedResults:       a.ExpectedResults,
		IsPayed:               a.IsPayed,
		AdditionalInformation: nil,
		Status:                domain.ApplicationStatus(a.Status),
		RejectedReason:        a.RejectedReason,
	}
}

func setApplicationDates(app *domain.Application, a applicationSeed) {
	if a.CreatedAt != nil {
		app.CreatedAt = addTimeShift(*a.CreatedAt, 24)
	} else {
		app.CreatedAt = time.Now()
	}
	if a.UpdatedAt != nil {
		app.UpdatedAt = addTimeShift(*a.UpdatedAt, 24)
	} else {
		app.UpdatedAt = time.Now()
	}
}

func addTimeShift(t time.Time, maxHours int) time.Time {
	shift := time.Duration(rand.Intn(maxHours*2)-maxHours) * time.Hour
	return t.Add(shift)
}
