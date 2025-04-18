package service

import "github.com/gin-gonic/gin"

func init() {
	registerGetHandler("/ping", handlePing)
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
