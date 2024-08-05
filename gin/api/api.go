package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/unifuu/hitotose/gin/model/game"
	"github.com/unifuu/hitotose/gin/srv"
	"github.com/unifuu/hitotose/gin/util"

	"github.com/gin-gonic/gin"
)

var (
	service srv.Service
)

const (
	PAGE_LIMIT = 20
)

func Init(e *gin.Engine) {
	service = srv.NewService()

	api := e.Group("/api/game")
	{
		api.GET("/", index)
		api.Any("/create", create)
		api.Any("/delete", delete)
		api.Any("/update", update)
		api.POST("/update_status", updateStatus)
		api.GET("/status", query)
		api.GET("/badges", badges)
		api.GET("/rank/target", target)
		api.GET("/status/:status/:platform/:page", status)
	}
}

func target(c *gin.Context) {

}

func updateStatus(c *gin.Context) {
	gameId := c.PostForm("id")
	newStatus := c.PostForm("newStatus")
	targetGame := service.ByID(gameId)
	targetGame.Status = game.Status(newStatus)
	service.Update(targetGame)
}

func badges(c *gin.Context) {
	status := game.Status(c.Query("status"))
	badges := service.Badge(status)
	c.JSON(http.StatusOK, gin.H{
		"badges": badges,
	})
}

func create(c *gin.Context) {
	switch c.Request.Method {
	case "POST":
		title := c.PostForm("title")
		developer := c.PostForm("developer")
		publisher := c.PostForm("publisher")
		genre := c.PostForm("genre")
		platform := c.PostForm("platform")

		service.Create(game.Game{
			Title:     title,
			Status:    game.PLAYING,
			Genre:     genre,
			Platform:  platform,
			Developer: developer,
			Publisher: publisher,
		})
		c.Redirect(http.StatusSeeOther, "/game")
	}
}

func delete(c *gin.Context) {
	id := c.Query("id")

	// Delete from db.game
	err := service.Delete(id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusNotFound, err)
		return
	}

	// Delete from assets
	file := id + ".webp"
	root := util.Root()
	path := root + "/../assets/img/games/" + file
	os.Remove(path)
}

func index(c *gin.Context) {
	status := game.PLAYING
	data := gin.H{
		"games": service.ByStatus(status),
	}
	c.JSON(http.StatusOK, data)
}

func query(c *gin.Context) {
	status := game.Status(c.Query("status"))
	data := gin.H{
		"games": service.ByStatus(status),
	}
	c.JSON(http.StatusOK, data)
}

func update(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		id := c.Query("id")
		g := service.ByID(id)

		c.JSON(http.StatusOK, gin.H{
			"game": g,
		})

	case "POST":
		gId := c.PostForm("id")
		developer := c.PostForm("developer")
		publisher := c.PostForm("publisher")

		playedHour, _ := strconv.Atoi(c.PostForm("played_hour"))
		playedMin, _ := strconv.Atoi(c.PostForm("played_min"))
		playedTime := playedHour*60 + playedMin

		hltbHour, _ := strconv.Atoi(c.PostForm("hltb_hour"))
		hltbMin, _ := strconv.Atoi(c.PostForm("hltb_min"))
		hltbTime := hltbHour*60 + hltbMin

		ranking, _ := strconv.Atoi(c.PostForm("ranking"))
		// rating, _ := strconv.Atoi(c.PostForm("rating"))

		g := service.ByID(gId)
		g.Title = c.PostForm("title")
		g.Developer = developer
		g.Publisher = publisher
		g.Status = game.Status(c.PostForm("status"))
		g.PlayedTime = playedTime
		g.TimeToBeat = hltbTime
		g.Genre = c.PostForm("genre")
		g.Platform = c.PostForm("platform")
		g.Ranking = ranking
		g.Rating = c.PostForm("rating")

		file, err := c.FormFile("cover")
		// Upload image to assets
		if file != nil && err == nil {
			fn := gId + ".webp"
			root := util.Root()
			path := root + "/assets/images/games/" + fn
			c.SaveUploadedFile(file, path)
		}
		service.Update(g)
		c.Redirect(http.StatusSeeOther, "/game")
	}
}

func status(c *gin.Context) {
	status := game.Status(c.Param("status"))
	platform := game.Platform(c.Param("platform"))
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		page = 1
	}

	games, totalPage := service.PageByPlatformStatus(status, platform, page, PAGE_LIMIT)
	data := gin.H{
		"games":      games,
		"total_page": totalPage,
	}

	c.JSON(http.StatusOK, data)
}
