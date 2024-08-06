package main

import (
	"github.com/gin-gonic/gin"
	"github.com/unifuu/hitotose/gin/api"
	"github.com/unifuu/hitotose/gin/db"
)

func main() {
	db.Init()

	// Create a new Gin router
	router := gin.Default()

	api.Init(router)

	// Run the server on port 8080
	router.Run(":8080")
}
