package controllers

import (
	"net/http"
	"strconv"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	permissions := c.GetStringSlice(helpers.CtxValues.Permissions)
	if !helpers.SliceStringContains(permissions, helpers.UserPermissions.ReadUser) {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(helpers.Unauthorized))
		c.Abort()
		return
	}

	var limit, pageNo int64
	var err error

	if c.Query("limit") != "" {
		limit, err = strconv.ParseInt(c.Query("limit"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse("Invalid Limit"))
			c.Abort()
			return
		}
	}
	if c.Query("pageNo") != "" {
		pageNo, err = strconv.ParseInt(c.Query("pageNo"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse("Invalid Page No."))
			c.Abort()
			return
		}
	}

	result, err := services.DBConn.FindUsers(c, nil, limit, pageNo)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("", result))
}
