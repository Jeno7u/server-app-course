package handlers

import (
	"net/http"
	"strconv"

	"github.com/Jeno7u/server-app-course/internal/models"
	"github.com/Jeno7u/server-app-course/internal/storage"
	"github.com/gin-gonic/gin"
)

func fetchCurrentUser(c *gin.Context) (models.CurrentUser, bool) {
	val, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Unauthorized"})
		return models.CurrentUser{}, false
	}
	return val.(models.CurrentUser), true
}

func CreateTask(c *gin.Context) {
	user, ok := fetchCurrentUser(c)
	if !ok {
		return
	}

	var req models.TaskCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		OwnerID:     user.ID,
	}

	created := storage.GlobalTaskStorage.Create(task)
	c.JSON(http.StatusCreated, created)
}

func GetTasks(c *gin.Context) {
	user, ok := fetchCurrentUser(c)
	if !ok {
		return
	}

	statusFilter := c.Query("status")
	minPriorityStr := c.Query("min_priority")
	minPriority := 0
	if minPriorityStr != "" {
		if val, err := strconv.Atoi(minPriorityStr); err == nil {
			minPriority = val
		}
	}

	allTasks := storage.GlobalTaskStorage.GetAll()
	filtered := make([]models.Task, 0)

	for _, t := range allTasks {
		if t.OwnerID != user.ID {
			continue // Only own tasks
		}
		if statusFilter != "" && t.Status != statusFilter {
			continue
		}
		if minPriority > 0 && t.Priority < minPriority {
			continue
		}
		filtered = append(filtered, t)
	}

	c.JSON(http.StatusOK, filtered)
}

func GetTaskByID(c *gin.Context) {
	user, ok := fetchCurrentUser(c)
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	task, exists := storage.GlobalTaskStorage.GetByID(id)
	if !exists || task.OwnerID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func UpdateTaskStatus(c *gin.Context) {
	user, ok := fetchCurrentUser(c)
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	var req models.TaskStatusUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	task, exists := storage.GlobalTaskStorage.GetByID(id)
	if !exists || task.OwnerID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	task.Status = req.Status
	storage.GlobalTaskStorage.Update(task)

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	user, ok := fetchCurrentUser(c)
	if !ok {
		return
	}

	id, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	task, exists := storage.GlobalTaskStorage.GetByID(id)
	if !exists || task.OwnerID != user.ID {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	storage.GlobalTaskStorage.Delete(id)
	c.Status(http.StatusNoContent)
}
