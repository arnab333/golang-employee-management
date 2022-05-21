package controllers

import (
	"net/http"
	"strconv"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetHolidays(c *gin.Context) {
	permissions := c.GetStringSlice(helpers.CtxValues.Permissions)
	if !helpers.SliceStringContains(permissions, helpers.UserPermissions.ReadHoliday) {
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

	result, err := services.DBConn.FindHolidays(c, nil, limit, pageNo)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("", result))
}

func GetHoliday(c *gin.Context) {
	permissions := c.GetStringSlice(helpers.CtxValues.Permissions)
	if !helpers.SliceStringContains(permissions, helpers.UserPermissions.ReadHoliday) {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(helpers.Unauthorized))
		c.Abort()
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusUnprocessableEntity, helpers.HandleErrorResponse(helpers.RequiredID))
		c.Abort()
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.HandleErrorResponse(helpers.InvalidID))
		c.Abort()
		return
	}

	filters := bson.M{
		"_id": objID,
	}

	result, err := services.DBConn.FindHoliday(c, filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("", result))
}

func CreateHoliday(c *gin.Context) {

}
