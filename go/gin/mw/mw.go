package mw

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	go_redis "github.com/go-redis/redis/v8"
	"github.com/unifuu/hitotose/gin/db/redis"
)

// Check user authority
func Auth(c *gin.Context) {
	authToken, _ := c.Cookie("auth_token")

	if authToken == "undefined" {
		authToken = ""
	}

	// Get auth token from the header
	// if len(authToken) == 0 {
	// 	if len(c.Request.Header["auth_token"]) > 0 {
	// 		authToken = c.Request.Header["auth_token"][0]
	// 	}
	// }

	uid, err := redis.Get(authToken)
	switch {
	case err == go_redis.Nil:
	case err != nil:
		fmt.Println(err.Error())
	}

	if uid == "" {
		c.Abort()
		return
	}
	c.Next()
}

// Redirect HTTP to HTTPS
func HTTPS() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.TLS == nil && c.Request.Host != "" {
			url := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusMovedPermanently, url)
			c.Abort()
			return
		}

		// If the request is using HTTPS then c.Next()
		c.Next()
	}
}
