package middleware

import (
	"github.com/gin-gonic/gin"
)

// MustAuth ensures authentication for Gin
func MustAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Cookie("auth"); err != nil || cookie == "" {
			c.Redirect(302, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}
