package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	UserID          uuid.UUID  `db:"user_id" json:"user_id"`
	TemplateID      uuid.UUID  `db:"template_id" json:"template_id"`
	Title           string     `db:"title" json:"title"`
	Slug            string     `db:"slug" json:"slug"`
	EventDate       time.Time  `db:"event_date" json:"event_date"`
	LocationName    *string    `db:"location_name" json:"location_name"`
	LocationAddress *string    `db:"location_address" json:"location_address"`
	IsPublished     bool       `db:"is_published" json:"is_published"`
	ViewCount       int        `db:"view_count" json:"view_count"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`

	// Relations (populated on demand)
	Theme    *EventTheme    `db:"-" json:"theme,omitempty"`
	Sections []EventSection `db:"-" json:"sections,omitempty"`
}

type EventTheme struct {
	ID             uuid.UUID `db:"id" json:"id"`
	EventID        uuid.UUID `db:"event_id" json:"event_id"`
	PrimaryColor   *string   `db:"primary_color" json:"primary_color"`
	SecondaryColor *string   `db:"secondary_color" json:"secondary_color"`
	FontFamily     *string   `db:"font_family" json:"font_family"`
	BackgroundURL  *string   `db:"background_url" json:"background_url"`
	CustomCSS      *string   `db:"custom_css" json:"custom_css"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
}

type EventSection struct {
	ID                uuid.UUID       `db:"id" json:"id"`
	EventID           uuid.UUID       `db:"event_id" json:"event_id"`
	TemplateSectionID uuid.UUID       `db:"template_section_id" json:"template_section_id"`
	Content           json.RawMessage `db:"content" json:"content"`
	IsVisible         bool            `db:"is_visible" json:"is_visible"`
	SortOrder         int             `db:"sort_order" json:"sort_order"`
}

// Request / Response types
type CreateEventRequest struct {
	TemplateID      string  `json:"template_id" binding:"required,uuid"`
	Title           string  `json:"title" binding:"required,min=3,max=200"`
	EventDate       string  `json:"event_date" binding:"required"`
	LocationName    *string `json:"location_name"`
	LocationAddress *string `json:"location_address"`
}

type UpdateEventRequest struct {
	Title           *string `json:"title"`
	EventDate       *string `json:"event_date"`
	LocationName    *string `json:"location_name"`
	LocationAddress *string `json:"location_address"`
}

type UpdateThemeRequest struct {
	PrimaryColor   *string `json:"primary_color"`
	SecondaryColor *string `json:"secondary_color"`
	FontFamily     *string `json:"font_family"`
	BackgroundURL  *string `json:"background_url"`
	CustomCSS      *string `json:"custom_css"`
}

type UpdateSectionRequest struct {
	Content   json.RawMessage `json:"content"`
	IsVisible *bool           `json:"is_visible"`
	SortOrder *int            `json:"sort_order"`
}

type PublicEventResponse struct {
	Event    *Event         `json:"event"`
	Theme    *EventTheme    `json:"theme"`
	Sections []EventSection `json:"sections"`
	Gallery  []Media        `json:"gallery"`
	Stats    *EventStats    `json:"stats"`
}

type EventStats struct {
	TotalRSVP      int `json:"total_rsvp"`
	TotalAttending int `json:"total_attending"`
	TotalDeclined  int `json:"total_declined"`
	TotalPending   int `json:"total_pending"`
	TotalMessages  int `json:"total_messages"`
}

type EventRepository interface {
	Create(ctx context.Context, event *Event) error
	FindByID(ctx context.Context, id uuid.UUID) (*Event, error)
	FindBySlug(ctx context.Context, slug string) (*Event, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]Event, error)
	Update(ctx context.Context, event *Event) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
	SlugExists(ctx context.Context, slug string) (bool, error)

	// Theme
	UpsertTheme(ctx context.Context, theme *EventTheme) error
	FindThemeByEventID(ctx context.Context, eventID uuid.UUID) (*EventTheme, error)

	// Sections
	CreateSections(ctx context.Context, sections []EventSection) error
	FindSectionsByEventID(ctx context.Context, eventID uuid.UUID) ([]EventSection, error)
	UpdateSection(ctx context.Context, section *EventSection) error

	// Stats
	GetStats(ctx context.Context, eventID uuid.UUID) (*EventStats, error)
}
