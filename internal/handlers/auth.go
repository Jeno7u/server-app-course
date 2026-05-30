package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Jeno7u/server-app-course/internal/auth"
	"github.com/Jeno7u/server-app-course/internal/db"
	"github.com/Jeno7u/server-app-course/internal/models"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username=?)", req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"detail": "User already exists"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Error hashing password"})
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", req.Username, hashedPassword, "user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "New user created"})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	var storedHash, role, storedUsername string
	err := db.DB.QueryRow("SELECT username, password, role FROM users WHERE username=?", req.Username).Scan(&storedUsername, &storedHash, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"detail": "Database error"})
		}
		return
	}

	// Double check username to prevent timing attacks slightly
	if !auth.ConstantTimeCompare(req.Username, storedUsername) {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Authorization failed"})
		return
	}

	if !auth.CheckPasswordHash(req.Password, storedHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"detail": "Authorization failed"})
		return
	}

	token, err := auth.CreateToken(req.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"token_type":   "bearer",
	})
}

func ProtectedResource(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Access granted"})
}

func AdminResource(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome Admin!"})
}

func UserResource(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome User/Admin!"})
}
