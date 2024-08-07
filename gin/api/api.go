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
	sw      *game.StopWatch
)

const (
	PAGE_LIMIT = 20
)

func Init(e *gin.Engine) {
	service = srv.NewService()

	api := e.Group("/api/game")
	{
		api.GET("/", index)
		api.GET("/badge", badges)
		api.POST("/create", create)
		api.Any("/delete", delete)
		api.GET("/pages", pages)
		api.GET("/start", start)
		api.GET("/stop", stop)
		api.GET("/stopwatch", stopwatch)
		api.GET("/status", query)
		api.GET("/terminate", terminate)
		api.Any("/update", update)
	}
}

func badges(c *gin.Context) {
	status := game.Status(c.Query("status"))
	badges := service.Badge(status)
	c.JSON(http.StatusOK, gin.H{
		"badges": badges,
	})
}

func create(c *gin.Context) {
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

		playedHour, _ := strconv.Atoi(c.PostForm("played_time_hour"))
		playedMin, _ := strconv.Atoi(c.PostForm("played_time_min"))
		playedTime := playedHour*60 + playedMin

		toBeatHour, _ := strconv.Atoi(c.PostForm("time_to_beat_hour"))
		toBeatMin, _ := strconv.Atoi(c.PostForm("time_to_beat_min"))
		toBeatTime := toBeatHour*60 + toBeatMin

		ranking, _ := strconv.Atoi(c.PostForm("ranking"))
		// rating, _ := strconv.Atoi(c.PostForm("rating"))

		g := service.ByID(gId)
		g.Title = c.PostForm("title")
		g.Developer = developer
		g.Publisher = publisher
		g.Status = game.Status(c.PostForm("status"))
		g.PlayedTime = playedTime
		g.TimeToBeat = toBeatTime
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

func pages(c *gin.Context) {
	status := game.Status(c.Query("status"))
	platform := game.Platform(c.Query("platform"))
	page, err := strconv.Atoi(c.Query("page"))
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

func start(c *gin.Context) {
	if sw != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "stopwatch is running",
		})
		return
	}

	id := c.Query("id")
	title := service.ByID(id).Title

	sw = game.NewStopWatch(id, title)
	err := sw.Start()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "bad request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "stopwatch started",
	})
}

func stop(c *gin.Context) {
	dur := sw.Stop()

	if dur > 0 {
		g := service.ByID(sw.GameID)
		g.PlayedTime += dur
		service.Update(g)
	}

	sw = nil
	c.Redirect(http.StatusSeeOther, "/game")
}

func stopwatch(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"stopwatch": sw,
	})
}

func terminate(c *gin.Context) {
	sw = nil
	c.JSON(http.StatusOK, nil)
}
