package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/galihaleanda/event-invitation/internal/domain"
)

type eventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) domain.EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	query := `
		INSERT INTO events (id, user_id, template_id, title, slug, event_date, location_name, location_address, is_published, view_count, created_at, updated_at)
		VALUES (:id, :user_id, :template_id, :title, :slug, :event_date, :location_name, :location_address, :is_published, :view_count, :created_at, :updated_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("eventRepository.Create: %w", err)
	}
	return nil
}

func (r *eventRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	var event domain.Event
	query := `SELECT * FROM events WHERE id = $1`
	if err := r.db.GetContext(ctx, &event, query, id); err != nil {
		return nil, fmt.Errorf("eventRepository.FindByID: %w", err)
	}
	return &event, nil
}

func (r *eventRepository) FindBySlug(ctx context.Context, slug string) (*domain.Event, error) {
	var event domain.Event
	query := `SELECT * FROM events WHERE slug = $1`
	if err := r.db.GetContext(ctx, &event, query, slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("eventRepository.FindBySlug: %w", err)
	}
	return &event, nil
}

func (r *eventRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Event, error) {
	var events []domain.Event
	query := `SELECT * FROM events WHERE user_id = $1 ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &events, query, userID); err != nil {
		return nil, fmt.Errorf("eventRepository.FindByUserID: %w", err)
	}
	return events, nil
}

func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	query := `
		UPDATE events SET
			title = :title,
			event_date = :event_date,
			location_name = :location_name,
			location_address = :location_address,
			is_published = :is_published,
			updated_at = :updated_at
		WHERE id = :id AND user_id = :user_id
	`
	res, err := r.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("eventRepository.Update: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("event not found or unauthorized")
	}
	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("eventRepository.Delete: %w", err)
	}
	return nil
}

func (r *eventRepository) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `UPDATE events SET view_count = view_count + 1 WHERE id = $1`, id)
	return err
}

func (r *eventRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `SELECT COUNT(1) FROM events WHERE slug = $1`, slug)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Theme

func (r *eventRepository) UpsertTheme(ctx context.Context, theme *domain.EventTheme) error {
	query := `
		INSERT INTO event_themes (id, event_id, primary_color, secondary_color, font_family, background_url, custom_css, created_at)
		VALUES (:id, :event_id, :primary_color, :secondary_color, :font_family, :background_url, :custom_css, :created_at)
		ON CONFLICT (event_id) DO UPDATE SET
			primary_color = EXCLUDED.primary_color,
			secondary_color = EXCLUDED.secondary_color,
			font_family = EXCLUDED.font_family,
			background_url = EXCLUDED.background_url,
			custom_css = EXCLUDED.custom_css
	`
	_, err := r.db.NamedExecContext(ctx, query, theme)
	if err != nil {
		return fmt.Errorf("eventRepository.UpsertTheme: %w", err)
	}
	return nil
}

func (r *eventRepository) FindThemeByEventID(ctx context.Context, eventID uuid.UUID) (*domain.EventTheme, error) {
	var theme domain.EventTheme
	query := `SELECT * FROM event_themes WHERE event_id = $1`
	if err := r.db.GetContext(ctx, &theme, query, eventID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("eventRepository.FindThemeByEventID: %w", err)
	}
	return &theme, nil
}

// Sections

func (r *eventRepository) CreateSections(ctx context.Context, sections []domain.EventSection) error {
	if len(sections) == 0 {
		return nil
	}
	query := `
		INSERT INTO event_sections (id, event_id, template_section_id, content, is_visible, sort_order)
		VALUES (:id, :event_id, :template_section_id, :content, :is_visible, :sort_order)
	`
	_, err := r.db.NamedExecContext(ctx, query, sections)
	if err != nil {
		return fmt.Errorf("eventRepository.CreateSections: %w", err)
	}
	return nil
}

func (r *eventRepository) FindSectionsByEventID(ctx context.Context, eventID uuid.UUID) ([]domain.EventSection, error) {
	var sections []domain.EventSection
	query := `SELECT * FROM event_sections WHERE event_id = $1 ORDER BY sort_order ASC`
	if err := r.db.SelectContext(ctx, &sections, query, eventID); err != nil {
		return nil, fmt.Errorf("eventRepository.FindSectionsByEventID: %w", err)
	}
	return sections, nil
}

func (r *eventRepository) UpdateSection(ctx context.Context, section *domain.EventSection) error {
	query := `
		UPDATE event_sections SET
			content = :content,
			is_visible = :is_visible,
			sort_order = :sort_order
		WHERE id = :id AND event_id = :event_id
	`
	_, err := r.db.NamedExecContext(ctx, query, section)
	if err != nil {
		return fmt.Errorf("eventRepository.UpdateSection: %w", err)
	}
	return nil
}

// Stats

func (r *eventRepository) GetStats(ctx context.Context, eventID uuid.UUID) (*domain.EventStats, error) {
	var stats domain.EventStats
	query := `
		SELECT
			COUNT(*) as total_rsvp,
			COUNT(*) FILTER (WHERE rsvp_status = 'yes') as total_attending,
			COUNT(*) FILTER (WHERE rsvp_status = 'no') as total_declined,
			COUNT(*) FILTER (WHERE rsvp_status = 'pending') as total_pending,
			COUNT(*) FILTER (WHERE message IS NOT NULL AND message != '') as total_messages
		FROM guests
		WHERE event_id = $1
	`
	if err := r.db.GetContext(ctx, &stats, query, eventID); err != nil {
		return nil, fmt.Errorf("eventRepository.GetStats: %w", err)
	}
	return &stats, nil
}
