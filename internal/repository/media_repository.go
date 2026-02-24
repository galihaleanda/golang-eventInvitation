package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/galihaleanda/event-invitation/internal/domain"
)

type mediaRepository struct {
	db *sqlx.DB
}

func NewMediaRepository(db *sqlx.DB) domain.MediaRepository {
	return &mediaRepository{db: db}
}

func (r *mediaRepository) Create(ctx context.Context, media *domain.Media) error {
	query := `
		INSERT INTO media (id, event_id, file_url, media_type, created_at)
		VALUES (:id, :event_id, :file_url, :media_type, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, media)
	if err != nil {
		return fmt.Errorf("mediaRepository.Create: %w", err)
	}
	return nil
}

func (r *mediaRepository) FindByEventID(ctx context.Context, eventID uuid.UUID) ([]domain.Media, error) {
	var media []domain.Media
	query := `SELECT * FROM media WHERE event_id = $1 ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &media, query, eventID); err != nil {
		return nil, fmt.Errorf("mediaRepository.FindByEventID: %w", err)
	}
	return media, nil
}

func (r *mediaRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM media WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("mediaRepository.Delete: %w", err)
	}
	return nil
}
