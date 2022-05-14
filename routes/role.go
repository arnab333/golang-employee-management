package routes

import (
	"github.com/arnab333/golang-employee-management/controllers"
	"github.com/arnab333/golang-employee-management/middlewares"
	"github.com/gin-gonic/gin"
)

func roleRoutes(rg *gin.RouterGroup) {
	rg.GET("/roles", middlewares.VerifyToken, middlewares.VerifyRole, controllers.GetUserRoles)

	rg.PUT("/roles", middlewares.VerifyToken, middlewares.VerifyRole, controllers.UpdateUserRole)
}
