package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/Jeno7u/server-app-course/internal/db"
	"github.com/Jeno7u/server-app-course/internal/models"
	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context) {
	var req models.TodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	completed := false
	if req.Completed != nil {
		completed = *req.Completed
	}

	res, err := db.DB.Exec("INSERT INTO todos (title, description, completed) VALUES (?, ?, ?)", req.Title, req.Description, completed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create todo"})
		return
	}

	id, _ := res.LastInsertId()

	todo := models.Todo{
		ID:          int(id),
		Title:       req.Title,
		Description: req.Description,
		Completed:   completed,
	}

	c.JSON(http.StatusCreated, todo)
}

func GetTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID"})
		return
	}

	var todo models.Todo
	err = db.DB.QueryRow("SELECT id, title, description, completed FROM todos WHERE id = ?", id).
		Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"detail": "Todo not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, todo)
}

func UpdateTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID"})
		return
	}

	var req models.TodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var exists bool
	_ = db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM todos WHERE id=?)", id).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Todo not found"})
		return
	}

	completed := false
	if req.Completed != nil {
		completed = *req.Completed
	}

	_, err = db.DB.Exec("UPDATE todos SET title = ?, description = ?, completed = ? WHERE id = ?", req.Title, req.Description, completed, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to update todo"})
		return
	}

	todo := models.Todo{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Completed:   completed,
	}

	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID"})
		return
	}

	res, err := db.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Database error"})
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"detail": "Successfully deleted"})
}
