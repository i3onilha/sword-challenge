package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"sword-challenge/internal/models"
	"sword-challenge/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService *service.TaskService
}

func NewTaskController(taskService *service.TaskService) *TaskController {
	return &TaskController{
		taskService: taskService,
	}
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required" example:"Fix air conditioning"`
	Summary     string `json:"summary" binding:"required" example:"Replaced filters and recharged coolant"`
	PerformedAt string `json:"performed_at" binding:"required" example:"2024-03-20T14:30:00Z"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"required" example:"Fix air conditioning"`
	Summary     string `json:"summary" binding:"required" example:"Replaced filters and recharged coolant"`
	PerformedAt string `json:"performed_at" binding:"required" example:"2024-03-20T14:30:00Z"`
}

// @Summary      Create a new task
// @Description  Create a new task for the authenticated technician
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task body CreateTaskRequest true "Task Information"
// @Success      201  {object}  models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/tasks [post]
func (h *TaskController) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	performedAt, err := parseTime(req.PerformedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid performed_at date format"})
		return
	}

	task := &models.Task{
		Title:       req.Title,
		Summary:     req.Summary,
		PerformedAt: performedAt,
	}

	userID := getUserIDFromContext(c)
	if err := h.taskService.CreateTask(c.Request.Context(), task, userID); err != nil {
		switch err {
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case service.ErrInvalidInput:
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid input"})
		default:
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, task)
}

// @Summary      Get a specific task
// @Description  Get a task by its ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path int true "Task ID"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/tasks/{id} [get]
func (h *TaskController) GetTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userID := getUserIDFromContext(c)
	task, err := h.taskService.GetTask(c.Request.Context(), taskID, userID)
	if err != nil {
		switch err {
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary      Get all tasks
// @Description  Get all tasks for the authenticated user
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Task
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/tasks [get]
func (h *TaskController) GetTasks(c *gin.Context) {
	userID := getUserIDFromContext(c)
	tasks, err := h.taskService.GetTasks(c.Request.Context(), userID)
	if err != nil {
		switch err {
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary      Update a task
// @Description  Update an existing task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path int true "Task ID"
// @Param        task body UpdateTaskRequest true "Task Information"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/tasks/{id} [put]
func (h *TaskController) UpdateTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	performedAt, err := parseTime(req.PerformedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid performed_at date format"})
		return
	}

	task := &models.Task{
		ID:          taskID,
		Summary:     req.Summary,
		PerformedAt: performedAt,
	}

	userID := getUserIDFromContext(c)
	if err := h.taskService.UpdateTask(c.Request.Context(), task, userID); err != nil {
		switch err {
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary      Delete a task
// @Description  Delete a task by its ID
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path int true "Task ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/tasks/{id} [delete]
func (h *TaskController) DeleteTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userID := getUserIDFromContext(c)
	if err := h.taskService.DeleteTask(c.Request.Context(), taskID, userID); err != nil {
		switch err {
		case service.ErrUnauthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		case service.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
