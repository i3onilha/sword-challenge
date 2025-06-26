package controllers

import (
	"net/http"
	"strconv"

	_ "sword-challenge/internal/models"
	"sword-challenge/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	notificationService *service.NotificationService
}

func NewNotificationController(notificationService *service.NotificationService) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
	}
}

// @Summary      Get unread notifications
// @Description  Get all unread notifications for the authenticated manager
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Notification
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/notifications [get]
func (h *NotificationController) GetUnreadNotifications(c *gin.Context) {
	userID := getUserIDFromContext(c)
	notifications, err := h.notificationService.GetUnreadNotifications(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// @Summary      Mark notification as read
// @Description  Mark a specific notification as read
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Param        id path int true "Notification ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/notifications/{id}/read [put]
func (h *NotificationController) MarkAsRead(c *gin.Context) {
	notificationID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	userID := getUserIDFromContext(c)
	if err := h.notificationService.MarkAsRead(c.Request.Context(), notificationID, userID); err != nil {
		switch err {
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
