package routes

import (
	"github.com/arnab333/golang-employee-management/controllers"
	"github.com/arnab333/golang-employee-management/middlewares"
	"github.com/gin-gonic/gin"
)

func authRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", controllers.Register)

	rg.POST("/login", controllers.Login)

	rg.POST("/logout", middlewares.VerifyToken, controllers.Logout)
}
