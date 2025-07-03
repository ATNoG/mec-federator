package middleware

import (
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func TimingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/federation/v1/ewbi") {
			start := time.Now()
			
			c.Next()
			
			duration := time.Since(start)
			log.Printf("EWBI Request %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
		} else {
			c.Next()
		}
	}
}