package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Jeno7u/server-app-course/internal/auth"
	"github.com/gin-gonic/gin"
)

var rateStore = struct {
	sync.Mutex
	m map[string][]int64
}{m: make(map[string][]int64)}

func RateLimit(limit int, windowSec int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		path := c.FullPath()
		key := clientIP + ":" + path

		now := time.Now().Unix()

		rateStore.Lock()
		defer rateStore.Unlock()

		arr := rateStore.m[key]
		cutoff := now - int64(windowSec)
		var filtered []int64

		for _, t := range arr {
			if t > cutoff {
				filtered = append(filtered, t)
			}
		}

		if len(filtered) >= limit {
			rateStore.m[key] = filtered
			c.JSON(http.StatusTooManyRequests, gin.H{"detail": "Too many requests"})
			c.Abort()
			return
		}

		filtered = append(filtered, now)
		rateStore.m[key] = filtered

		c.Next()
	}
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"detail": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("username", claims["sub"])
		c.Set("role", claims["role"])
		c.Next()
	}
}

func RoleAuth(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Role not found"})
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Invalid role type"})
			c.Abort()
			return
		}

		allowed := false
		for _, r := range allowedRoles {
			if r == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"detail": "Access forbidden: insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func DocsAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		mode := strings.ToUpper(os.Getenv("MODE"))
		if mode == "PROD" {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		if mode == "DEV" {
			expectedUser := os.Getenv("DOCS_USER")
			expectedPass := os.Getenv("DOCS_PASSWORD")

			user, pass, hasAuth := c.Request.BasicAuth()
			if !hasAuth || subtle.ConstantTimeCompare([]byte(user), []byte(expectedUser)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(expectedPass)) != 1 {
				c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
				c.JSON(http.StatusUnauthorized, gin.H{"detail": "Unauthorized"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
