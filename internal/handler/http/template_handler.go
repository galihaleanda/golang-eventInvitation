package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/galihaleanda/event-invitation/internal/service"
	"github.com/galihaleanda/event-invitation/internal/utils"
)

type TemplateHandler struct {
	templateService service.TemplateService
}

func NewTemplateHandler(templateService service.TemplateService) *TemplateHandler {
	return &TemplateHandler{templateService: templateService}
}

func (h *TemplateHandler) GetAll(c *gin.Context) {
	category := c.Query("category")
	templates, err := h.templateService.GetAll(c.Request.Context(), category)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "failed to get templates")
		return
	}
	utils.RespondOK(c, templates)
}

func (h *TemplateHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	tmpl, err := h.templateService.GetByID(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := err.(*service.AppError); ok {
			utils.RespondError(c, appErr.Code, appErr.Message)
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "failed to get template")
		return
	}
	utils.RespondOK(c, tmpl)
}
