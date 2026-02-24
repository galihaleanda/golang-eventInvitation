package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/galihaleanda/event-invitation/internal/domain"
	"github.com/galihaleanda/event-invitation/internal/service"
	"github.com/galihaleanda/event-invitation/internal/utils"
)

type RSVPHandler struct {
	rsvpService service.RSVPService
}

func NewRSVPHandler(rsvpService service.RSVPService) *RSVPHandler {
	return &RSVPHandler{rsvpService: rsvpService}
}

// POST /events/:id/rsvp  (public)
func (h *RSVPHandler) Submit(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid event id")
		return
	}

	var req domain.RSVPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	guest, err := h.rsvpService.Submit(c.Request.Context(), eventID, &req)
	if err != nil {
		if appErr, ok := err.(*service.AppError); ok {
			utils.RespondError(c, appErr.Code, appErr.Message)
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "failed to submit rsvp")
		return
	}
	utils.RespondCreated(c, guest)
}

// GET /events/:id/guests  (protected - owner only)
func (h *RSVPHandler) GetGuests(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid event id")
		return
	}

	guests, err := h.rsvpService.GetGuests(c.Request.Context(), getUserID(c), eventID)
	if err != nil {
		if appErr, ok := err.(*service.AppError); ok {
			utils.RespondError(c, appErr.Code, appErr.Message)
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "failed to get guests")
		return
	}
	utils.RespondOK(c, guests)
}
