package api

import (
	"github.com/adishjain1107/tradex/pkg/order-matcher/app"
	"github.com/gin-gonic/gin"
)

func Routes(application *app.App) *gin.Engine {
	_ = application
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "order-matcher ping pong",
		})
	})

	return router
}
