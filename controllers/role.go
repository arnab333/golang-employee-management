package controllers

import (
	"net/http"
	"strings"

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

	result, err := services.DBConn.FindRoles(c, nil)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("", result))
}

func UpdateUserRole(c *gin.Context) {
	var json services.UserRole
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		return
	}

	if strings.Contains(json.ID.Hex(), "000000000000000000000000") {
		c.JSON(http.StatusUnprocessableEntity, helpers.HandleErrorResponse(helpers.RequiredID))
		c.Abort()
		return
	}

	objID, err := primitive.ObjectIDFromHex(json.ID.Hex())
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

	result, err := services.DBConn.UpdateRole(c, filters, update)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		return
	}

	msg := "Role Updated!"
	if result.UpsertedID != nil {
		msg = "Role Inserted!"
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("", msg))
}
