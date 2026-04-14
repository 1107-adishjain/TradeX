package controller

import (
	"errors"
	"net/http"

	"github.com/adishjain1107/tradex/pkg/auth/app"
	"github.com/adishjain1107/tradex/pkg/auth/models"
	"github.com/adishjain1107/tradex/pkg/auth/service"
	"github.com/gin-gonic/gin"
)

func Register(application *app.App) gin.HandlerFunc {
	authservice := service.NewAuthService(application.DB)
	return func(c *gin.Context) {
		var req models.RegReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := authservice.RegisterService(c.Request.Context(), req)

		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidPayload):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, service.ErrEmailUppercase):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, service.ErrUserAlreadyExists):
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "registered successfully",
			"user":    user,
		})

	}
}

func Login(application *app.App) gin.HandlerFunc {
	authservice := service.NewAuthService(application.DB)

	return func(c *gin.Context) {
		var req models.LoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		authResp, err := authservice.LoginService(c.Request.Context(), req)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidPayload):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, service.ErrEmailUppercase):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.Is(err, service.ErrInvalidCredentials):
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
			return
		}

		c.SetCookie(
			"refresh_token",
			authResp.RefreshToken,
			7*24*60*60,
			"/",
			"",
			true,
			true,
		)

		c.JSON(http.StatusOK, gin.H{
			"message":      "login successful",
			"access_token": authResp.AccessToken,
			"email":        authResp.Email,
			"role":         authResp.Role,
		})
	}
}
