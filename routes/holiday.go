package routes

import (
	"github.com/arnab333/golang-employee-management/controllers"
	"github.com/gin-gonic/gin"
)

func holidayRoutes(rg *gin.RouterGroup) {
	rg.GET("/holidays", controllers.GetHolidays)

	rg.GET("/holiday/:id", controllers.GetHoliday)

}
