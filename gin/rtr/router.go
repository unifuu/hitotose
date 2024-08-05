package router

import (
	"github.com/gin-gonic/gin"
	"github.com/unifuu/hitotose/gin/api/game"
)

func Route(e *gin.Engine) {
	game.Init(e)
}
