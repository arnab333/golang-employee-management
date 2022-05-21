package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserRoles(c *gin.Context) {
	permissions := c.GetStringSlice(helpers.CtxValues.Permissions)
	if !helpers.SliceStringContains(permissions, helpers.UserPermissions.ReadRole) {
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

	result, err := services.DBConn.FindRoles(c, nil, limit, pageNo)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("", result))
}

func GetUserRole(c *gin.Context) {
	permissions := c.GetStringSlice(helpers.CtxValues.Permissions)
	if !helpers.SliceStringContains(permissions, helpers.UserPermissions.ReadRole) {
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

	filters := bson.M{"_id": objID}

	result, err := services.DBConn.FindRole(c, filters)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("", result))
}

func CreateUserRole(c *gin.Context) {
	permissions := c.GetStringSlice(helpers.CtxValues.Permissions)
	if !helpers.SliceStringContains(permissions, helpers.UserPermissions.CreateRole) {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(helpers.Unauthorized))
		c.Abort()
		return
	}

	var json services.UserRole
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	if json.Name == "" {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse("`name` is required"))
		c.Abort()
		return
	}

	filters := bson.M{"name": json.Name}
	result, _ := services.DBConn.FindRole(c, filters)
	if result.Name != "" {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse("Role already exists!"))
		c.Abort()
		return
	}

	if _, err := services.DBConn.InsertRole(c, json); err != nil {
		fmt.Println("InsertRole err", err.Error())
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("Role Inserted!", nil))
}

func UpdateUserRole(c *gin.Context) {
	permissions := c.GetStringSlice(helpers.CtxValues.Permissions)
	if !helpers.SliceStringContains(permissions, helpers.UserPermissions.UpdateRole) {
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

	var json services.UserRole
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.HandleErrorResponse(helpers.InvalidID))
		c.Abort()
		return
	}

	filters := bson.M{"_id": objID}

	update := bson.M{
		"$set": bson.M{
			"permissions": json.Permissions,
		},
	}

	_, err = services.DBConn.UpdateRole(c, filters, update)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("Role Updated!", nil))
}
