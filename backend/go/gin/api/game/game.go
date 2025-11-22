package game

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/unifuu/hitotose/backend/go/gin/model/game"
	"github.com/unifuu/hitotose/backend/go/gin/mw"
	game_svc "github.com/unifuu/hitotose/backend/go/gin/svc/game"
	"github.com/unifuu/hitotose/backend/go/gin/util"

	"github.com/gin-gonic/gin"
)

const (
	PAGE_LIMIT = 15
)

var (
	svc game_svc.Service
	sw  *game.StopWatch
)

func Init(e *gin.Engine) {
	svc = game_svc.NewService()

	anon := e.Group("/api/game")
	{
		anon.GET("/", index)
		anon.GET("/badge", badges)
		anon.GET("/pages", pages)
	}

	auth := e.Group("/api/game").Use(mw.Auth)
	{
		auth.POST("/create", create)
		auth.Any("/delete", delete)
		auth.GET("/start", start)
		auth.GET("/stop", stop)
		auth.GET("/stopwatch", stopwatch)
		auth.GET("/status", query)
		auth.GET("/terminate", terminate)
		auth.Any("/update", update)
		auth.POST("/update/rating", updateRating)
	}
}

func badges(c *gin.Context) {
	status := game.Status(c.Query("status"))
	badges := svc.Badge(status)
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

	svc.Create(game.Game{
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
	err := svc.Delete(id)
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
		"games": svc.ByStatus(status),
	}
	c.JSON(http.StatusOK, data)
}

func query(c *gin.Context) {
	status := game.Status(c.Query("status"))
	data := gin.H{
		"games": svc.ByStatus(status),
	}
	c.JSON(http.StatusOK, data)
}

func update(c *gin.Context) {
	switch c.Request.Method {
	case "GET":
		id := c.Query("id")
		g := svc.ByID(id)

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

		g := svc.ByID(gId)
		g.Title = c.PostForm("title")
		g.Developer = developer
		g.Publisher = publisher
		g.Status = game.Status(c.PostForm("status"))
		g.PlayedTime = playedTime
		g.Genre = c.PostForm("genre")
		g.Platform = c.PostForm("platform")
		rating, _ := strconv.Atoi(c.PostForm("rating"))
		g.Rating = rating

		file, err := c.FormFile("cover")
		// Upload image to assets
		if file != nil && err == nil {
			fn := gId + ".webp"
			root := util.Root()
			path := root + "/assets/images/games/" + fn
			c.SaveUploadedFile(file, path)
		}
		err = svc.Update(g)
		if err != nil {
			log.Println(err)
		}
		c.Redirect(http.StatusSeeOther, "/game")
	}
}

func updateRating(c *gin.Context) {
	gId := c.PostForm("id")
	g := svc.ByID(gId)
	rating, _ := strconv.Atoi(c.PostForm("rating"))
	g.Rating = rating

	err := svc.Update(g)
	if err != nil {
		log.Println(err)
	}
	c.Redirect(http.StatusSeeOther, "/game")
}

func pages(c *gin.Context) {
	keyword := c.Query("keyword")
	status := game.Status(c.Query("status"))
	platform := game.Platform(c.Query("platform"))
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	games, totalPage := svc.Query(keyword, platform, status, page, PAGE_LIMIT)
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
	title := svc.ByID(id).Title

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
		g := svc.ByID(sw.GameID)
		g.PlayedTime += dur
		svc.Update(g)
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
