package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/galihaleanda/event-invitation/internal/domain"
)

type templateRepository struct {
	db *sqlx.DB
}

func NewTemplateRepository(db *sqlx.DB) domain.TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) FindAll(ctx context.Context, category string) ([]domain.Template, error) {
	var templates []domain.Template
	query := `SELECT * FROM templates WHERE is_active = true`
	args := []interface{}{}

	if category != "" {
		query += ` AND category = $1`
		args = append(args, category)
	}
	query += ` ORDER BY created_at DESC`

	if err := r.db.SelectContext(ctx, &templates, query, args...); err != nil {
		return nil, fmt.Errorf("templateRepository.FindAll: %w", err)
	}
	return templates, nil
}

func (r *templateRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Template, error) {
	var tmpl domain.Template
	query := `SELECT * FROM templates WHERE id = $1`
	if err := r.db.GetContext(ctx, &tmpl, query, id); err != nil {
		return nil, fmt.Errorf("templateRepository.FindByID: %w", err)
	}
	return &tmpl, nil
}

func (r *templateRepository) FindSectionsByTemplateID(ctx context.Context, templateID uuid.UUID) ([]domain.TemplateSection, error) {
	var sections []domain.TemplateSection
	query := `SELECT * FROM template_sections WHERE template_id = $1 ORDER BY sort_order ASC`
	if err := r.db.SelectContext(ctx, &sections, query, templateID); err != nil {
		return nil, fmt.Errorf("templateRepository.FindSectionsByTemplateID: %w", err)
	}
	return sections, nil
}
