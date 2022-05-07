package controllers

import (
	"fmt"
	"net/http"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var json services.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		return
	}

	json.IsActive = true

	if json.Role != helpers.UserRoles.Admin && json.Role != helpers.UserRoles.Superadmin && json.Role != helpers.UserRoles.User {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse("Role Does Not Match!"))
		return
	}

	result, err := services.DBConn.FindUser(c, bson.M{"email": json.Email})

	if err != nil {
		fmt.Println("FindUser err ==>", err.Error())
	}

	if result.Email != "" {
		c.JSON(http.StatusConflict, helpers.HandleErrorResponse("Email Already Exists!"))
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(json.Password), 14)

	json.Password = string(password)

	if err != nil {
		fmt.Println("FindUser err ==>", err.Error())
	}

	if _, err := services.DBConn.InsertUser(c, json); err != nil {
		fmt.Println("InsertUser err ==>", err)
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helpers.HandleSuccessResponse(helpers.CreatedMessage, nil))
}

func Login(c *gin.Context) {

}

func Logout(c *gin.Context) {

}
