package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RSVPStatus string

const (
	RSVPStatusPending RSVPStatus = "pending"
	RSVPStatusYes     RSVPStatus = "yes"
	RSVPStatusNo      RSVPStatus = "no"
)

type Guest struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	EventID    uuid.UUID  `db:"event_id" json:"event_id"`
	Name       string     `db:"name" json:"name"`
	Phone      *string    `db:"phone" json:"phone"`
	Message    *string    `db:"message" json:"message"`
	RSVPStatus RSVPStatus `db:"rsvp_status" json:"rsvp_status"`
	GuestCode  *string    `db:"guest_code" json:"guest_code"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
}

type RSVPRequest struct {
	Name      string     `json:"name" binding:"required,min=2,max=150"`
	Phone     *string    `json:"phone"`
	Message   *string    `json:"message"`
	Status    RSVPStatus `json:"status" binding:"required,oneof=yes no pending"`
	GuestCode *string    `json:"guest_code"`
}

type GuestRepository interface {
	Create(ctx context.Context, guest *Guest) error
	FindByEventID(ctx context.Context, eventID uuid.UUID) ([]Guest, error)
	FindByGuestCode(ctx context.Context, code string) (*Guest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status RSVPStatus) error
}
