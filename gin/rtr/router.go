package router

import (
	"github.com/gin-gonic/gin"
	"github.com/unifuu/hitotose/gin/api"
)

func Route(e *gin.Engine) {
	api.Init(e)
}
