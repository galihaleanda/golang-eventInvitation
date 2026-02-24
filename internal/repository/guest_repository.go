package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/galihaleanda/event-invitation/internal/domain"
)

type guestRepository struct {
	db *sqlx.DB
}

func NewGuestRepository(db *sqlx.DB) domain.GuestRepository {
	return &guestRepository{db: db}
}

func (r *guestRepository) Create(ctx context.Context, guest *domain.Guest) error {
	query := `
		INSERT INTO guests (id, event_id, name, phone, message, rsvp_status, guest_code, created_at)
		VALUES (:id, :event_id, :name, :phone, :message, :rsvp_status, :guest_code, :created_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, guest)
	if err != nil {
		return fmt.Errorf("guestRepository.Create: %w", err)
	}
	return nil
}

func (r *guestRepository) FindByEventID(ctx context.Context, eventID uuid.UUID) ([]domain.Guest, error) {
	var guests []domain.Guest
	query := `SELECT * FROM guests WHERE event_id = $1 ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &guests, query, eventID); err != nil {
		return nil, fmt.Errorf("guestRepository.FindByEventID: %w", err)
	}
	return guests, nil
}

func (r *guestRepository) FindByGuestCode(ctx context.Context, code string) (*domain.Guest, error) {
	var guest domain.Guest
	query := `SELECT * FROM guests WHERE guest_code = $1`
	if err := r.db.GetContext(ctx, &guest, query, code); err != nil {
		return nil, fmt.Errorf("guestRepository.FindByGuestCode: %w", err)
	}
	return &guest, nil
}

func (r *guestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RSVPStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE guests SET rsvp_status = $1 WHERE id = $2`,
		status, id,
	)
	if err != nil {
		return fmt.Errorf("guestRepository.UpdateStatus: %w", err)
	}
	return nil
}
