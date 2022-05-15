package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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

	role, err := services.DBConn.FindRole(c, bson.M{"name": json.Role})

	if err != nil {
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

	if err != nil {
		fmt.Println("FindUser err ==>", err.Error())
	}

	json.Password = string(password)
	json.Role = role.ID.Hex()

	if _, err := services.DBConn.InsertUser(c, json); err != nil {
		fmt.Println("InsertUser err ==>", err)
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		return
	}

	mailResponse, err := services.SendEmail(&services.EmailDetails{
		Name:    json.FullName,
		Address: json.Email,
	})

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("StatusCode ==>", mailResponse.StatusCode)
	}

	c.JSON(http.StatusCreated, helpers.HandleSuccessResponse(helpers.CreatedMessage, nil))
}

func Login(c *gin.Context) {
	var json map[string]string

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		return
	}

	user, err := services.DBConn.FindUser(c, bson.M{"email": json["email"], "isActive": true})

	if err != nil {
		c.JSON(http.StatusNotFound, helpers.HandleErrorResponse("Invalid Email or Password!!"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json["password"])); err != nil {
		c.JSON(http.StatusNotFound, helpers.HandleErrorResponse("Invalid Email or Password!"))
		return
	}

	td, err := services.CreateAuth(c, user.ID.Hex(), user.Role)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	tokens := gin.H{
		"accessToken":  td.AccessToken,
		"refreshToken": td.RefreshToken,
	}
	c.JSON(http.StatusCreated, helpers.HandleSuccessResponse("", tokens))
}

func Logout(c *gin.Context) {
	accessUUID := c.GetString(helpers.CtxValues.AccessUUID)
	if accessUUID == "" {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Unauthorized!"))
		return
	}

	deleted, err := services.DeleteAuth(c, accessUUID)
	if err != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Unauthorized!!"))
		return
	}
	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("Successfully logged out", nil))
}

func RefreshToken(c *gin.Context) {
	var json map[string]string

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.HandleErrorResponse(err.Error()))
		return
	}

	tokenString := json["refreshToken"]

	claims, err := services.ExtractFromToken(tokenString, helpers.EnvKeys.JWT_REFRESH_SECRET)
	if err != nil {
		msg := "Invalid Token!"
		if strings.Contains(err.Error(), "expired") {
			msg = "Token Expired!"
		}
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse(msg))
		c.Abort()
		return
	}

	deleted, err := services.DeleteAuth(c, claims.ID)
	if err != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Unauthorized!!"))
		return
	}

	td, err := services.CreateAuth(c, claims.UserID, claims.Role)
	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	tokens := gin.H{
		"accessToken":  td.AccessToken,
		"refreshToken": td.RefreshToken,
	}
	c.JSON(http.StatusCreated, helpers.HandleSuccessResponse("", tokens))
}
