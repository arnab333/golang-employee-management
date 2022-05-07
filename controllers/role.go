package controllers

import (
	"net/http"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {
	var data []string = []string{
		helpers.UserRoles.User,
		helpers.UserRoles.Admin,
		helpers.UserRoles.Superadmin,
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("", data))
}
