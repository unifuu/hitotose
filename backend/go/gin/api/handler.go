package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/unifuu/hitotose/gin/api/game"
	"github.com/unifuu/hitotose/gin/api/user"
)

func Init(e *gin.Engine) {
	assets(e)

	game.Init(e)
	user.Init(e)
}

func assets(e *gin.Engine) {
	e.Static("/assets", "./assets")
}
