package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()

	config.AllowOriginFunc = func(origin string) bool {
		return true
	}

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	return cors.New(config) //New(config) returns a CORS middleware with the user-defined configuration. It will panic if the configuration is invalid.

}
