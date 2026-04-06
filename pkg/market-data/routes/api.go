package api

import (
	mw "github.com/adishjain1107/tradex/pkg/common/middleware"
	"github.com/adishjain1107/tradex/pkg/market-data/app"
	"github.com/gin-gonic/gin"
)

func Routes(application *app.App) *gin.Engine {
	_ = application
	router := gin.Default()

	router.Use(mw.SecurityHeaders())
	router.Use(mw.CorsMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "market-data ping pong",
		})
	})

	return router
}
