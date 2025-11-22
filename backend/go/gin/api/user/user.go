package user

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/unifuu/hitotose/backend/go/gin/db/redis"
	"github.com/unifuu/hitotose/backend/go/gin/model/user"

	"github.com/gin-gonic/gin"
	user_svc "github.com/unifuu/hitotose/backend/go/gin/svc/user"
)

const (
	// 3 months
	REDIS_EXPIRATION = time.Duration(24 * time.Hour * 30 * 3)

	// 3 months
	GIN_COOKIE_EXPIRATION = 1 * 60 * 60 * 24 * 30 * 3
)

var svc user_svc.Service

func Init(e *gin.Engine) {
	svc = user_svc.NewService()

	e.Any("/api/user/checkAuth", checkAuth)
	e.Any("/api/user/checkToken", checkToken)
}

// Handle user checkAuth
func checkAuth(c *gin.Context) {
	var u user.User
	c.BindJSON(&u)

	// Check the username and password is valid or not
	auth, err := svc.SignIn(u.Username, u.Password)

	if err != nil {
		respJson(c, false, "")
	} else {
		// Generate a token and set it to Redis
		token := setAuth(c, auth.ID.Hex())

		respJson(c, true, token)
	}
}

// checkToken checks the auth token is expired or not
func checkToken(c *gin.Context) {
	var req struct {
		AuthToken string `json:"auth_token"`
	}

	// Get auth token from cookie
	authToken, _ := c.Cookie("auth_token")

	// Get auth token from request header
	if len(authToken) == 0 {
		authToken = c.Request.Header.Get("auth_token")
	}

	err := c.BindJSON(&req)
	if err == nil {
		authToken = req.AuthToken
	}

	// Get auth token from url
	if len(authToken) == 0 {
		authToken = c.Query("auth_token")
	}

	// Cannot get auth token
	if len(authToken) == 0 {
		respJson(c, false, "")
		return
	}

	// Get user id from Redis
	userID, _ := redis.Get(authToken)

	// Check the user id is exist or not
	if len(userID) > 0 {
		respJson(c, true, authToken)
	} else {
		respJson(c, false, "")
	}
}

func clearAuth(c *gin.Context) {
	authToken, _ := c.Cookie("auth_token")
	redis.Del(authToken)
	c.SetCookie("auth_token", "", -1, "/", "", false, false)
}

func respJson(c *gin.Context, isAuth bool, token string) {
	if isAuth {
		c.JSON(http.StatusOK, gin.H{
			"is_auth":    true,
			"auth_token": token,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"is_auth":    false,
			"auth_token": token,
		})
	}
}

// Set user authentication token to Redis
func setAuth(c *gin.Context, userID string) string {
	b := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("Failed to generate a random value...")
	}
	token := base64.URLEncoding.EncodeToString(b)

	if err := redis.Set(token, userID, REDIS_EXPIRATION); err != nil {
		panic("Failed to set session key to Redis..." + err.Error())
	}

	c.SetCookie("auth_token", token, GIN_COOKIE_EXPIRATION, "/", "", false, false)
	return token
}
