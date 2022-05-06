package routes

import (
	"github.com/arnab333/golang-employee-management/controllers"
	"github.com/gin-gonic/gin"
)

func roleRoutes(rg *gin.RouterGroup) {
	rg.POST("/role", controllers.CreateRole)
}
