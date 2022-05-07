package routes

import (
	"github.com/arnab333/golang-employee-management/controllers"
	"github.com/gin-gonic/gin"
)

func authRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", controllers.Register)

	rg.POST("/login", controllers.Login)

	rg.POST("/logout", controllers.Logout)
}
