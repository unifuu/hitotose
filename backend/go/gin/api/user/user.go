package user

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/unifuu/hitotose/backend/go/gin/db/redis"
	"github.com/unifuu/hitotose/backend/go/gin/model/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	user_svc "github.com/unifuu/hitotose/backend/go/gin/svc/user"
)

const (
	REDIS_EXPIRATION      = time.Duration(24 * time.Hour * 30 * 3)
	GIN_COOKIE_EXPIRATION = 1 * 60 * 60 * 24 * 30 * 3
)

var (
	svc       user_svc.Service
	jwtSecret string // JWT secret key
	isSecure  bool   // for HTTPS-only cookies
)

func Init(e *gin.Engine) {
	svc = user_svc.NewService()

	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		panic("JWT_SECRET environment variable is not set")
	}

	e.Any("/api/user/checkAuth", checkAuth)
	e.Any("/api/user/checkToken", checkToken)
	e.POST("/api/user/logout", logout)
}

// Handle user checkAuth
func checkAuth(c *gin.Context) {
	var u user.User

	// Bind and validate JSON
	err := c.BindJSON(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"is_auth": false,
			"error":   "Invalid request body",
		})
		return
	}

	// Validate input
	if u.Username == "" || u.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"is_auth": false,
			"error":   "Username and password are required",
		})
		return
	}

	// Check the username and password is valid or not
	auth, err := svc.SignIn(u.Username, u.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"is_auth": false,
			"error":   "Invalid username or password",
		})
	} else {
		// Generate a token and set it to Redis
		token, err := setAuth(c, auth.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"is_auth": false,
				"error":   "Failed to create session",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"is_auth":    true,
			"auth_token": token,
		})
	}
}

// checkToken checks the auth token
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
	if err == nil && req.AuthToken != "" {
		authToken = req.AuthToken
	}

	// Cannot get auth token
	if len(authToken) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"is_auth": false,
			"error":   "No authentication token provided",
		})
		return
	}

	// Verify JWT signature and expiration
	userID, err := verifyJWT(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"is_auth": false,
			"error":   "Invalid or expired token",
		})
		return
	}

	// Get user id from Redis
	storedUserID, err := redis.Get(authToken)

	// Check the user id exists and matches
	if err != nil || len(storedUserID) == 0 || storedUserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"is_auth": false,
			"error":   "Session not found or expired",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"is_auth":    true,
			"auth_token": authToken,
		})
	}
}

// logout handles user logout
func logout(c *gin.Context) {
	authToken, err := c.Cookie("auth_token")
	if err != nil || authToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No active session",
		})
		return
	}

	// Delete from Redis
	if err := redis.Del(authToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout",
		})
		return
	}

	// Clear cookie
	c.SetCookie("auth_token", "", -1, "/", "", isSecure, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// Set user authentication token to Redis
func setAuth(c *gin.Context, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(REDIS_EXPIRATION).Unix(),
	})

	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	// Save session token to Redis
	err = redis.Set(signed, userID, REDIS_EXPIRATION)
	if err != nil {
		return "", err
	}

	// Set secure cookie
	// httpOnly=true prevents JavaScript access (XSS protection)
	// secure=true ensures HTTPS only in production
	c.SetCookie(
		"auth_token",          // name
		signed,                // value
		GIN_COOKIE_EXPIRATION, // maxAge in seconds
		"/",                   // path
		"",                    // domain
		isSecure,              // secure (HTTPS only in production)
		true,                  // httpOnly (prevent JavaScript access)
	)

	return signed, nil
}

// verifyJWT validates the JWT token and returns the user ID
func verifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	// Extract user ID from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("user ID not found in token")
	}

	return userID, nil
}
