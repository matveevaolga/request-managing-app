package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matveevaolga/request-managing-app/internal/domain"
	domainrepo "github.com/matveevaolga/request-managing-app/internal/domain/repository"
)

type applicationRepository struct {
	db      *pgxpool.Pool
	builder squirrel.StatementBuilderType
}

func NewApplicationRepository(db *pgxpool.Pool) *applicationRepository {
	return &applicationRepository{
		db:      db,
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *applicationRepository) Create(ctx context.Context, app *domain.Application) error {
	query := r.builder.Insert("applications").
		Columns("full_name", "email", "phone", "organisation_name", "organisation_url",
			"project_name", "type_id", "expected_results", "is_payed", "additional_information",
			"status", "created_at", "updated_at").
		Values(app.FullName, app.Email, app.Phone, app.OrganisationName, app.OrganisationURL,
			app.ProjectName, app.TypeID, app.ExpectedResults, app.IsPayed, app.AdditionalInformation,
			app.Status, app.CreatedAt, app.UpdatedAt).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&app.ID)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}
	return nil
}

func (r *applicationRepository) GetByID(ctx context.Context, id int64) (*domain.Application, error) {
	query := r.builder.Select(
		"id", "full_name", "email", "phone", "organisation_name", "organisation_url",
		"project_name", "type_id", "expected_results", "is_payed", "additional_information",
		"status", "rejected_reason", "reviewer", "created_at", "updated_at",
	).From("applications").Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var app domain.Application
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&app.ID, &app.FullName, &app.Email, &app.Phone, &app.OrganisationName, &app.OrganisationURL,
		&app.ProjectName, &app.TypeID, &app.ExpectedResults, &app.IsPayed, &app.AdditionalInformation,
		&app.Status, &app.RejectedReason, &app.ReviewerID, &app.CreatedAt, &app.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application by id: %w", err)
	}
	return &app, nil
}

func (r *applicationRepository) GetAllFiltered(ctx context.Context, params domainrepo.ApplicationFilterParameters) ([]domain.ApplicationPreview, int, error) {
	query := r.builder.Select(
		"a.id", "a.project_name", "pt.name as type_name", "a.full_name",
		"a.organisation_name", "a.updated_at", "a.status", "a.rejected_reason",
	).From("applications a").Join("project_types pt ON a.type_id = pt.id")

	countQuery := r.builder.Select("COUNT(*)").From("applications a")

	query, countQuery = r.filterApplications(query, countQuery, params)
	query = r.sortApplications(query, params.SortByDateUpdated)
	query = r.paginateApplications(query, params.Limit, params.Offset)

	count, countArgs, err := countQuery.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var total int
	err = r.db.QueryRow(ctx, count, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count applications: %w", err)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get filtered applications: %w", err)
	}
	defer rows.Close()

	applications, err := r.scanApplications(rows)
	if err != nil {
		return nil, 0, err
	}

	return applications, total, nil
}

func (r *applicationRepository) Update(ctx context.Context, app *domain.Application) error {
	query := r.builder.Update("applications").
		Set("full_name", app.FullName).
		Set("email", app.Email).
		Set("phone", app.Phone).
		Set("organisation_name", app.OrganisationName).
		Set("organisation_url", app.OrganisationURL).
		Set("project_name", app.ProjectName).
		Set("type_id", app.TypeID).
		Set("expected_results", app.ExpectedResults).
		Set("is_payed", app.IsPayed).
		Set("additional_information", app.AdditionalInformation).
		Set("status", app.Status).
		Set("rejected_reason", app.RejectedReason).
		Set("reviewer", app.ReviewerID).
		Set("updated_at", app.UpdatedAt).
		Where("id = ?", app.ID)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	cmdTag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrApplicationNotFound
	}
	return nil
}

func (r *applicationRepository) UpdateStatus(ctx context.Context, id int64, status domain.ApplicationStatus, reviewerID *int64, rejectedReason *string) error {
	query := r.builder.Update("applications").
		Set("status", status).
		Set("reviewer", reviewerID).
		Set("rejected_reason", rejectedReason).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where("id = ?", id)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	cmdTag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrApplicationNotFound
	}
	return nil
}

func (r *applicationRepository) filterApplications(query, countQuery squirrel.SelectBuilder, params domainrepo.ApplicationFilterParameters) (squirrel.SelectBuilder, squirrel.SelectBuilder) {
	if params.Active != nil {
		if *params.Active {
			query = query.Where("a.status = ?", domain.StatusPending)
			countQuery = countQuery.Where("status = ?", domain.StatusPending)
		} else {
			query = query.Where("a.status != ?", domain.StatusPending)
			countQuery = countQuery.Where("status != ?", domain.StatusPending)
		}
	}

	if params.Search != nil && *params.Search != "" {
		arg := "%" + *params.Search + "%"
		cond := squirrel.Or{
			squirrel.Expr("a.project_name ILIKE ?", arg),
			squirrel.Expr("a.full_name ILIKE ?", arg),
		}
		query = query.Where(cond)
		countQuery = countQuery.Where(cond)
	}

	if params.ProjectTypeID != nil {
		query = query.Where("a.type_id = ?", *params.ProjectTypeID)
		countQuery = countQuery.Where("type_id = ?", *params.ProjectTypeID)
	}

	return query, countQuery
}

func (r *applicationRepository) sortApplications(query squirrel.SelectBuilder, sortBy string) squirrel.SelectBuilder {
	switch sortBy {
	case "ASC":
		return query.OrderBy("a.updated_at ASC")
	case "DESC":
		return query.OrderBy("a.updated_at DESC")
	default:
		return query.OrderBy("a.updated_at DESC")
	}
}

func (r *applicationRepository) paginateApplications(query squirrel.SelectBuilder, limit, offset int) squirrel.SelectBuilder {
	if limit > 0 {
		query = query.Limit(uint64(limit))
		if offset > 0 {
			query = query.Offset(uint64(offset))
		}
	}
	return query
}

func (r *applicationRepository) scanApplications(rows pgx.Rows) ([]domain.ApplicationPreview, error) {
	var apps []domain.ApplicationPreview
	for rows.Next() {
		var app domain.ApplicationPreview
		err := rows.Scan(
			&app.ID, &app.ProjectName, &app.TypeName, &app.Initiator,
			&app.OrganisationName, &app.DateUpdated, &app.Status, &app.RejectionMessage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan application: %w", err)
		}
		apps = append(apps, app)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}
	return apps, nil
}
