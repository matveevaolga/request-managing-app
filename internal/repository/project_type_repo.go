package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matveevaolga/request-managing-app/internal/domain"
)

type projectTypeRepository struct {
	db *pgxpool.Pool
}

func NewProjectTypeRepository(db *pgxpool.Pool) *projectTypeRepository {
	return &projectTypeRepository{db: db}
}

func (r *projectTypeRepository) Create(ctx context.Context, pt *domain.ProjectType) error {
	query := `INSERT INTO project_types (name, created_at) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(ctx, query, pt.Name, pt.CreatedAt).Scan(&pt.ID)
	if err != nil {
		return fmt.Errorf("failed to create project type: %w", err)
	}
	return nil
}

func (r *projectTypeRepository) GetByID(ctx context.Context, id int64) (*domain.ProjectType, error) {
	var pt domain.ProjectType
	query := `SELECT id, name, created_at FROM project_types WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&pt.ID, &pt.Name, &pt.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProjectTypeNotFound
		}
		return nil, fmt.Errorf("failed to get project type by id: %w", err)
	}
	return &pt, nil
}

func (r *projectTypeRepository) GetByName(ctx context.Context, name string) (*domain.ProjectType, error) {
	var pt domain.ProjectType
	query := `SELECT id, name, created_at FROM project_types WHERE name = $1`
	err := r.db.QueryRow(ctx, query, name).Scan(&pt.ID, &pt.Name, &pt.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProjectTypeNotFound
		}
		return nil, fmt.Errorf("failed to get project type by name: %w", err)
	}
	return &pt, nil
}

func (r *projectTypeRepository) GetAllProjects(ctx context.Context) ([]domain.ProjectType, error) {
	query := `SELECT id, name, created_at FROM project_types ORDER BY id`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get project types: %w", err)
	}
	defer rows.Close()

	var pts []domain.ProjectType
	for rows.Next() {
		var pt domain.ProjectType
		err := rows.Scan(&pt.ID, &pt.Name, &pt.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project type: %w", err)
		}
		pts = append(pts, pt)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}
	return pts, nil
}

func (r *projectTypeRepository) Update(ctx context.Context, pt *domain.ProjectType) error {
	query := `UPDATE project_types SET name = $1 WHERE id = $2`
	cmdTag, err := r.db.Exec(ctx, query, pt.Name, pt.ID)
	if err != nil {
		return fmt.Errorf("failed to update project type: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrProjectTypeNotFound
	}
	return nil
}

func (r *projectTypeRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM project_types WHERE id = $1`
	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project type: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrProjectTypeNotFound
	}
	return nil
}
