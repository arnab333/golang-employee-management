package controllers

import (
	"net/http"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
)

func GetHolidays(c *gin.Context) {
	permissions := c.GetStringSlice(helpers.CtxValues.Permissions)
	if !helpers.SliceStringContains(permissions, helpers.UserPermissions.ReadHoliday) {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(helpers.Unauthorized))
		c.Abort()
		return
	}

	result, err := services.DBConn.FindHolidays(c, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("", result))
}
