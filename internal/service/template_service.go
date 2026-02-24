package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/galihaleanda/event-invitation/internal/domain"
)

type TemplateService interface {
	GetAll(ctx context.Context, category string) ([]domain.Template, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error)
}

type templateService struct {
	templateRepo domain.TemplateRepository
}

func NewTemplateService(templateRepo domain.TemplateRepository) TemplateService {
	return &templateService{templateRepo: templateRepo}
}

func (s *templateService) GetAll(ctx context.Context, category string) ([]domain.Template, error) {
	templates, err := s.templateRepo.FindAll(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("templateService.GetAll: %w", err)
	}
	return templates, nil
}

func (s *templateService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error) {
	tmpl, err := s.templateRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("templateService.GetByID: %w", err)
	}

	sections, err := s.templateRepo.FindSectionsByTemplateID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("templateService.GetByID (sections): %w", err)
	}
	tmpl.Sections = sections

	return tmpl, nil
}
