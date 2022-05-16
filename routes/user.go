package routes

import (
	"github.com/arnab333/golang-employee-management/controllers"
	"github.com/gin-gonic/gin"
)

func userRoutes(rg *gin.RouterGroup) {
	rg.GET("/users", controllers.GetUsers)

	rg.GET("/user/:id", func(c *gin.Context) {})

	rg.POST("/user", func(c *gin.Context) {})

	rg.PUT("/user", func(c *gin.Context) {})

}
