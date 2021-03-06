package routes

import (
	"github.com/arnab333/golang-employee-management/controllers"
	"github.com/gin-gonic/gin"
)

func roleRoutes(rg *gin.RouterGroup) {
	rg.GET("/roles", controllers.GetUserRoles)

	rg.PUT("/role/:id", controllers.UpdateUserRole)

	rg.GET("/role/:id", controllers.GetUserRole)

	rg.POST("/role", controllers.CreateUserRole)
}
