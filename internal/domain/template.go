package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID           uuid.UUID        `db:"id" json:"id"`
	Name         string           `db:"name" json:"name"`
	Category     string           `db:"category" json:"category"`
	ThumbnailURL *string          `db:"thumbnail_url" json:"thumbnail_url"`
	IsActive     bool             `db:"is_active" json:"is_active"`
	CreatedAt    time.Time        `db:"created_at" json:"created_at"`
	Sections     []TemplateSection `db:"-" json:"sections,omitempty"`
}

type TemplateSection struct {
	ID             uuid.UUID       `db:"id" json:"id"`
	TemplateID     uuid.UUID       `db:"template_id" json:"template_id"`
	Name           string          `db:"name" json:"name"`
	Type           string          `db:"type" json:"type"`
	DefaultContent json.RawMessage `db:"default_content" json:"default_content"`
	SortOrder      int             `db:"sort_order" json:"sort_order"`
}

type TemplateRepository interface {
	FindAll(ctx context.Context, category string) ([]Template, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Template, error)
	FindSectionsByTemplateID(ctx context.Context, templateID uuid.UUID) ([]TemplateSection, error)
}
