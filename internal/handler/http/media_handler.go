package http

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/galihaleanda/event-invitation/internal/config"
	"github.com/galihaleanda/event-invitation/internal/domain"
	"github.com/galihaleanda/event-invitation/internal/utils"
)

type MediaHandler struct {
	mediaRepo   domain.MediaRepository
	eventRepo   domain.EventRepository
	storageCfg  config.StorageConfig
}

func NewMediaHandler(mediaRepo domain.MediaRepository, eventRepo domain.EventRepository, cfg *config.Config) *MediaHandler {
	return &MediaHandler{
		mediaRepo:  mediaRepo,
		eventRepo:  eventRepo,
		storageCfg: cfg.Storage,
	}
}

// POST /events/:id/media
func (h *MediaHandler) Upload(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid event id")
		return
	}

	// Verify ownership
	event, err := h.eventRepo.FindByID(c.Request.Context(), eventID)
	if err != nil || event.UserID != getUserID(c) {
		utils.RespondError(c, http.StatusForbidden, "forbidden")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	// Determine media type from extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	mediaType := domain.MediaTypeImage
	switch ext {
	case ".mp4", ".mov", ".avi":
		mediaType = domain.MediaTypeVideo
	case ".mp3", ".wav", ".ogg":
		mediaType = domain.MediaTypeAudio
	}

	// Save file
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	savePath := filepath.Join(h.storageCfg.BasePath, "events", eventID.String(), filename)

	if err := c.SaveUploadedFile(header, savePath); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to save file")
		return
	}

	fileURL := fmt.Sprintf("%s/events/%s/%s", h.storageCfg.BaseURL, eventID.String(), filename)

	media := &domain.Media{
		ID:        uuid.New(),
		EventID:   eventID,
		FileURL:   fileURL,
		MediaType: mediaType,
		CreatedAt: time.Now(),
	}

	if err := h.mediaRepo.Create(c.Request.Context(), media); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to save media")
		return
	}

	utils.RespondCreated(c, media)
}

// GET /events/:id/media
func (h *MediaHandler) GetByEvent(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid event id")
		return
	}

	media, err := h.mediaRepo.FindByEventID(c.Request.Context(), eventID)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get media")
		return
	}
	utils.RespondOK(c, media)
}

// DELETE /events/:id/media/:mediaId
func (h *MediaHandler) Delete(c *gin.Context) {
	mediaID, err := uuid.Parse(c.Param("mediaId"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid media id")
		return
	}

	if err := h.mediaRepo.Delete(c.Request.Context(), mediaID); err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to delete media")
		return
	}
	utils.RespondOK(c, nil)
}
