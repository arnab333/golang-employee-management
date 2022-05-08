package controllers

import (
	"net/http"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserRoles(c *gin.Context) {
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

	filter := bson.D{{Key: "name", Value: json.Name}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: json.Name},
		}},
	}

	result, err := services.DBConn.UpdateRole(c, filter, update)

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
