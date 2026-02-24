package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/galihaleanda/event-invitation/internal/domain"
	"github.com/galihaleanda/event-invitation/internal/service"
	"github.com/galihaleanda/event-invitation/internal/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*service.AppError); ok {
			utils.RespondError(c, appErr.Code, appErr.Message)
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, "registered successfully", resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*service.AppError); ok {
			utils.RespondError(c, appErr.Code, appErr.Message)
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondOK(c, resp)
}
