package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/galihaleanda/event-invitation/internal/domain"
	"github.com/galihaleanda/event-invitation/internal/utils"
)

type EventService interface {
	Create(ctx context.Context, userID uuid.UUID, req *domain.CreateEventRequest) (*domain.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error)
	GetBySlug(ctx context.Context, slug string) (*domain.PublicEventResponse, error)
	GetMyEvents(ctx context.Context, userID uuid.UUID) ([]domain.Event, error)
	Update(ctx context.Context, userID, eventID uuid.UUID, req *domain.UpdateEventRequest) (*domain.Event, error)
	Delete(ctx context.Context, userID, eventID uuid.UUID) error
	Publish(ctx context.Context, userID, eventID uuid.UUID, publish bool) error
	UpdateTheme(ctx context.Context, userID, eventID uuid.UUID, req *domain.UpdateThemeRequest) (*domain.EventTheme, error)
	UpdateSection(ctx context.Context, userID, eventID, sectionID uuid.UUID, req *domain.UpdateSectionRequest) (*domain.EventSection, error)
}

type eventService struct {
	eventRepo    domain.EventRepository
	templateRepo domain.TemplateRepository
	mediaRepo    domain.MediaRepository
}

func NewEventService(
	eventRepo domain.EventRepository,
	templateRepo domain.TemplateRepository,
	mediaRepo domain.MediaRepository,
) EventService {
	return &eventService{
		eventRepo:    eventRepo,
		templateRepo: templateRepo,
		mediaRepo:    mediaRepo,
	}
}

func (s *eventService) Create(ctx context.Context, userID uuid.UUID, req *domain.CreateEventRequest) (*domain.Event, error) {
	templateID, err := uuid.Parse(req.TemplateID)
	if err != nil {
		return nil, NewAppError(http.StatusBadRequest, "invalid template_id")
	}

	// Validate template exists
	tmpl, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return nil, NewAppError(http.StatusNotFound, "template not found")
	}

	eventDate, err := time.Parse(time.RFC3339, req.EventDate)
	if err != nil {
		return nil, NewAppError(http.StatusBadRequest, "invalid event_date format, use RFC3339")
	}

	// Generate unique slug
	slug := s.generateUniqueSlug(ctx, req.Title)

	now := time.Now()
	event := &domain.Event{
		ID:              uuid.New(),
		UserID:          userID,
		TemplateID:      templateID,
		Title:           req.Title,
		Slug:            slug,
		EventDate:       eventDate,
		LocationName:    req.LocationName,
		LocationAddress: req.LocationAddress,
		IsPublished:     false,
		ViewCount:       0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Copy template sections to event sections
	templateSections, err := s.templateRepo.FindSectionsByTemplateID(ctx, tmpl.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch template sections: %w", err)
	}

	var eventSections []domain.EventSection
	for _, ts := range templateSections {
		content := ts.DefaultContent
		if content == nil {
			content = json.RawMessage(`{}`)
		}
		eventSections = append(eventSections, domain.EventSection{
			ID:                uuid.New(),
			EventID:           event.ID,
			TemplateSectionID: ts.ID,
			Content:           content,
			IsVisible:         true,
			SortOrder:         ts.SortOrder,
		})
	}

	if len(eventSections) > 0 {
		if err := s.eventRepo.CreateSections(ctx, eventSections); err != nil {
			return nil, fmt.Errorf("failed to create event sections: %w", err)
		}
	}

	event.Sections = eventSections
	return event, nil
}

func (s *eventService) generateUniqueSlug(ctx context.Context, title string) string {
	for {
		slug := utils.GenerateSlug(title)
		exists, _ := s.eventRepo.SlugExists(ctx, slug)
		if !exists {
			return slug
		}
	}
}

func (s *eventService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	event, err := s.eventRepo.FindByID(ctx, id)
	if err != nil {
		return nil, NewAppError(http.StatusNotFound, "event not found")
	}
	return event, nil
}

func (s *eventService) GetBySlug(ctx context.Context, slug string) (*domain.PublicEventResponse, error) {
	event, err := s.eventRepo.FindBySlug(ctx, slug)
	if err != nil || event == nil {
		return nil, NewAppError(http.StatusNotFound, "event not found")
	}
	if !event.IsPublished {
		return nil, NewAppError(http.StatusNotFound, "event not found")
	}

	// Increment view count (fire and forget)
	go s.eventRepo.IncrementViewCount(context.Background(), event.ID)

	theme, _ := s.eventRepo.FindThemeByEventID(ctx, event.ID)
	sections, _ := s.eventRepo.FindSectionsByEventID(ctx, event.ID)
	gallery, _ := s.mediaRepo.FindByEventID(ctx, event.ID)
	stats, _ := s.eventRepo.GetStats(ctx, event.ID)

	return &domain.PublicEventResponse{
		Event:    event,
		Theme:    theme,
		Sections: sections,
		Gallery:  gallery,
		Stats:    stats,
	}, nil
}

func (s *eventService) GetMyEvents(ctx context.Context, userID uuid.UUID) ([]domain.Event, error) {
	events, err := s.eventRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	return events, nil
}

func (s *eventService) Update(ctx context.Context, userID, eventID uuid.UUID, req *domain.UpdateEventRequest) (*domain.Event, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, NewAppError(http.StatusNotFound, "event not found")
	}
	if event.UserID != userID {
		return nil, NewAppError(http.StatusForbidden, "forbidden")
	}

	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.EventDate != nil {
		t, err := time.Parse(time.RFC3339, *req.EventDate)
		if err != nil {
			return nil, NewAppError(http.StatusBadRequest, "invalid event_date format")
		}
		event.EventDate = t
	}
	if req.LocationName != nil {
		event.LocationName = req.LocationName
	}
	if req.LocationAddress != nil {
		event.LocationAddress = req.LocationAddress
	}
	event.UpdatedAt = time.Now()

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to update event: %w", err)
	}
	return event, nil
}

func (s *eventService) Delete(ctx context.Context, userID, eventID uuid.UUID) error {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return NewAppError(http.StatusNotFound, "event not found")
	}
	if event.UserID != userID {
		return NewAppError(http.StatusForbidden, "forbidden")
	}
	return s.eventRepo.Delete(ctx, eventID)
}

func (s *eventService) Publish(ctx context.Context, userID, eventID uuid.UUID, publish bool) error {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return NewAppError(http.StatusNotFound, "event not found")
	}
	if event.UserID != userID {
		return NewAppError(http.StatusForbidden, "forbidden")
	}
	event.IsPublished = publish
	event.UpdatedAt = time.Now()
	return s.eventRepo.Update(ctx, event)
}

func (s *eventService) UpdateTheme(ctx context.Context, userID, eventID uuid.UUID, req *domain.UpdateThemeRequest) (*domain.EventTheme, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, NewAppError(http.StatusNotFound, "event not found")
	}
	if event.UserID != userID {
		return nil, NewAppError(http.StatusForbidden, "forbidden")
	}

	theme := &domain.EventTheme{
		ID:             uuid.New(),
		EventID:        eventID,
		PrimaryColor:   req.PrimaryColor,
		SecondaryColor: req.SecondaryColor,
		FontFamily:     req.FontFamily,
		BackgroundURL:  req.BackgroundURL,
		CustomCSS:      req.CustomCSS,
		CreatedAt:      time.Now(),
	}

	if err := s.eventRepo.UpsertTheme(ctx, theme); err != nil {
		return nil, fmt.Errorf("failed to update theme: %w", err)
	}
	return theme, nil
}

func (s *eventService) UpdateSection(ctx context.Context, userID, eventID, sectionID uuid.UUID, req *domain.UpdateSectionRequest) (*domain.EventSection, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, NewAppError(http.StatusNotFound, "event not found")
	}
	if event.UserID != userID {
		return nil, NewAppError(http.StatusForbidden, "forbidden")
	}

	// Get current sections to find the target
	sections, err := s.eventRepo.FindSectionsByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to find sections: %w", err)
	}

	var target *domain.EventSection
	for i := range sections {
		if sections[i].ID == sectionID {
			target = &sections[i]
			break
		}
	}
	if target == nil {
		return nil, NewAppError(http.StatusNotFound, "section not found")
	}

	if req.Content != nil {
		target.Content = req.Content
	}
	if req.IsVisible != nil {
		target.IsVisible = *req.IsVisible
	}
	if req.SortOrder != nil {
		target.SortOrder = *req.SortOrder
	}

	if err := s.eventRepo.UpdateSection(ctx, target); err != nil {
		return nil, fmt.Errorf("failed to update section: %w", err)
	}
	return target, nil
}
