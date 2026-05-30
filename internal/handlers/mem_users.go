package handlers

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type MemUserIn struct {
	Username string `json:"username" binding:"required"`
	Age      int    `json:"age" binding:"required"`
}

type MemUserOut struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Age      int    `json:"age"`
}

var (
	memDB    = make(map[int]MemUserOut)
	memDBCnt = 1
	memMut   sync.Mutex
)

func ResetMemDB() {
	memMut.Lock()
	defer memMut.Unlock()
	memDB = make(map[int]MemUserOut)
	memDBCnt = 1
}

func CreateMemUser(c *gin.Context) {
	var in MemUserIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}

	memMut.Lock()
	defer memMut.Unlock()

	out := MemUserOut{
		ID:       memDBCnt,
		Username: in.Username,
		Age:      in.Age,
	}

	memDB[memDBCnt] = out
	memDBCnt++

	c.JSON(http.StatusCreated, out)
}

func GetMemUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	memMut.Lock()
	defer memMut.Unlock()

	user, ok := memDB[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func DeleteMemUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid ID format"})
		return
	}

	memMut.Lock()
	defer memMut.Unlock()

	_, ok := memDB[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"detail": "User not found"})
		return
	}

	delete(memDB, id)
	c.Status(http.StatusNoContent)
}
