package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func authRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", func(c *gin.Context) {})

	rg.POST("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "it's signup",
		})
	})
	rg.POST("/logout", func(c *gin.Context) {})
}
