package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unifuu/hitotose/gin/db"
)

func main() {
	db.Init()

	// Create a new Gin router
	router := gin.Default()

	// Define a GET route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Define a POST route
	router.POST("/post-example", func(c *gin.Context) {
		var json map[string]interface{}
		if err := c.BindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"received": json})
	})

	// Run the server on port 8080
	router.Run(":8080")
}
