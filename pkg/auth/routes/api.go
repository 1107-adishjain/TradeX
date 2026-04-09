package api

import (
	"github.com/adishjain1107/tradex/pkg/auth/app"
	"github.com/adishjain1107/tradex/pkg/auth/controller"
	mw "github.com/adishjain1107/tradex/pkg/common/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(application *app.App) *gin.Engine {
	router := gin.Default()

	router.Use(mw.SecurityHeaders())
	router.Use(mw.CorsMiddleware())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping pong",
		})
	})

	r := router.Group("/api/v1/auth")
	{
		r.POST("/register", controller.Register(application))
		r.POST("/login", controller.Login(application))
	}

	return router
}
