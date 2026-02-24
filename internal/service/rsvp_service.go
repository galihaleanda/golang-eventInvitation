package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/galihaleanda/event-invitation/internal/domain"
)

type RSVPService interface {
	Submit(ctx context.Context, eventID uuid.UUID, req *domain.RSVPRequest) (*domain.Guest, error)
	GetGuests(ctx context.Context, userID, eventID uuid.UUID) ([]domain.Guest, error)
}

type rsvpService struct {
	guestRepo domain.GuestRepository
	eventRepo domain.EventRepository
}

func NewRSVPService(guestRepo domain.GuestRepository, eventRepo domain.EventRepository) RSVPService {
	return &rsvpService{guestRepo: guestRepo, eventRepo: eventRepo}
}

func (s *rsvpService) Submit(ctx context.Context, eventID uuid.UUID, req *domain.RSVPRequest) (*domain.Guest, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil || event == nil {
		return nil, NewAppError(http.StatusNotFound, "event not found")
	}
	if !event.IsPublished {
		return nil, NewAppError(http.StatusBadRequest, "event is not published")
	}

	guest := &domain.Guest{
		ID:         uuid.New(),
		EventID:    eventID,
		Name:       req.Name,
		Phone:      req.Phone,
		Message:    req.Message,
		RSVPStatus: req.Status,
		GuestCode:  req.GuestCode,
		CreatedAt:  time.Now(),
	}

	if err := s.guestRepo.Create(ctx, guest); err != nil {
		return nil, fmt.Errorf("failed to save rsvp: %w", err)
	}
	return guest, nil
}

func (s *rsvpService) GetGuests(ctx context.Context, userID, eventID uuid.UUID) ([]domain.Guest, error) {
	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, NewAppError(http.StatusNotFound, "event not found")
	}
	if event.UserID != userID {
		return nil, NewAppError(http.StatusForbidden, "forbidden")
	}

	guests, err := s.guestRepo.FindByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guests: %w", err)
	}
	return guests, nil
}
