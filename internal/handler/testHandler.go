package handler

import "github.com/gin-gonic/gin"

func TestHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Test Handler",
	})
}
