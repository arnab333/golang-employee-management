package routes

import "github.com/gin-gonic/gin"

func userRoutes(rg *gin.RouterGroup) {
	rg.GET("/users", func(c *gin.Context) {})

	rg.GET("/user/:id", func(c *gin.Context) {})

	rg.POST("/user", func(c *gin.Context) {})

	rg.PUT("/user", func(c *gin.Context) {})

	rg.DELETE("/user", func(c *gin.Context) {})
}
