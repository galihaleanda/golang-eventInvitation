package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
)

type Media struct {
	ID        uuid.UUID `db:"id" json:"id"`
	EventID   uuid.UUID `db:"event_id" json:"event_id"`
	FileURL   string    `db:"file_url" json:"file_url"`
	MediaType MediaType `db:"media_type" json:"media_type"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type MediaRepository interface {
	Create(ctx context.Context, media *Media) error
	FindByEventID(ctx context.Context, eventID uuid.UUID) ([]Media, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
