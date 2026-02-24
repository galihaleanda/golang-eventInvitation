package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/galihaleanda/event-invitation/internal/domain"
	"github.com/galihaleanda/event-invitation/internal/middleware"
	"github.com/galihaleanda/event-invitation/internal/service"
	"github.com/galihaleanda/event-invitation/internal/utils"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

func getUserID(c *gin.Context) uuid.UUID {
	return c.MustGet(middleware.UserIDKey).(uuid.UUID)
}

// POST /events
func (h *EventHandler) Create(c *gin.Context) {
	var req domain.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	event, err := h.eventService.Create(c.Request.Context(), getUserID(c), &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.RespondCreated(c, event)
}

// GET /events (my events)
func (h *EventHandler) GetMyEvents(c *gin.Context) {
	events, err := h.eventService.GetMyEvents(c.Request.Context(), getUserID(c))
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get events")
		return
	}
	utils.RespondOK(c, events)
}

// GET /events/:id
func (h *EventHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	event, err := h.eventService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Only owner can view unpublished events
	if !event.IsPublished && event.UserID != getUserID(c) {
		utils.RespondError(c, http.StatusForbidden, "forbidden")
		return
	}

	utils.RespondOK(c, event)
}

// PATCH /events/:id
func (h *EventHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req domain.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	event, err := h.eventService.Update(c.Request.Context(), getUserID(c), id, &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.RespondOK(c, event)
}

// DELETE /events/:id
func (h *EventHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.eventService.Delete(c.Request.Context(), getUserID(c), id); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.RespondOK(c, nil)
}

// PATCH /events/:id/publish
func (h *EventHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Publish bool `json:"publish"`
	}
	_ = c.ShouldBindJSON(&body)

	if err := h.eventService.Publish(c.Request.Context(), getUserID(c), id, body.Publish); err != nil {
		handleServiceError(c, err)
		return
	}
	utils.RespondOK(c, gin.H{"is_published": body.Publish})
}

// PUT /events/:id/theme
func (h *EventHandler) UpdateTheme(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	var req domain.UpdateThemeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	theme, err := h.eventService.UpdateTheme(c.Request.Context(), getUserID(c), id, &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.RespondOK(c, theme)
}

// PATCH /events/:id/sections/:sectionId
func (h *EventHandler) UpdateSection(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid event id")
		return
	}
	sectionID, err := uuid.Parse(c.Param("sectionId"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid section id")
		return
	}

	var req domain.UpdateSectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	section, err := h.eventService.UpdateSection(c.Request.Context(), getUserID(c), eventID, sectionID, &req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.RespondOK(c, section)
}

// GET /e/:slug  (public)
func (h *EventHandler) GetPublic(c *gin.Context) {
	slug := c.Param("slug")
	resp, err := h.eventService.GetBySlug(c.Request.Context(), slug)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	utils.RespondOK(c, resp)
}

func handleServiceError(c *gin.Context, err error) {
	if appErr, ok := err.(*service.AppError); ok {
		utils.RespondError(c, appErr.Code, appErr.Message)
		return
	}
	utils.RespondError(c, http.StatusInternalServerError, "internal server error")
}
