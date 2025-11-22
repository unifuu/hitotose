package handler

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/unifuu/hitotose/backend/go/gin/api/game"
	"github.com/unifuu/hitotose/backend/go/gin/api/user"
)

func Init(e *gin.Engine) {
	assets(e)
	e.Use(static.Serve("/", static.LocalFile("../../../frontend/react/build", true)))

	game.Init(e)
	user.Init(e)
}

func assets(e *gin.Engine) {
	e.Static("/assets", "./assets")
}
