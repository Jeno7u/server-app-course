package handlers

import (
	"net/http"
	"strconv"

	"github.com/Jeno7u/server-app-course/internal/storage"
	"github.com/gin-gonic/gin"
)

func AdminStats(c *gin.Context) {
	allTasks := storage.GlobalTaskStorage.GetAll()

	byStatus := map[string]int{
		"todo":        0,
		"in_progress": 0,
		"done":        0,
	}

	for _, t := range allTasks {
		byStatus[t.Status]++
	}

	c.JSON(http.StatusOK, gin.H{
		"total_tasks": len(allTasks),
		"by_status":   byStatus,
	})
}

func AdminDeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	ok := storage.GlobalTaskStorage.Delete(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Task not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
