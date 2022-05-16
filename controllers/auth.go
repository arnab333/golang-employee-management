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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var json services.User

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	json.IsActive = true

	role, err := services.DBConn.FindRole(c, bson.M{"name": json.Role})

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse("Role Does Not Match!"))
		c.Abort()
		return
	}

	result, err := services.DBConn.FindUser(c, bson.M{"email": json.Email})

	if err != nil {
		fmt.Println("FindUser err ==>", err.Error())
	}

	if result.Email != "" {
		c.JSON(http.StatusConflict, helpers.HandleErrorResponse("Email Already Exists!"))
		c.Abort()
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(json.Password), 14)
	if err != nil {
		fmt.Println("GenerateFromPassword err ==>", err.Error())
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	json.Password = string(password)
	json.Role = role.ID.Hex()

	if _, err := services.DBConn.InsertUser(c, json); err != nil {
		fmt.Println("InsertUser err ==>", err)
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
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
		c.Abort()
		return
	}

	user, err := services.DBConn.FindUser(c, bson.M{"email": json["email"], "isActive": true})
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.HandleErrorResponse("Invalid Email or Password!!"))
		c.Abort()
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json["password"])); err != nil {
		c.JSON(http.StatusNotFound, helpers.HandleErrorResponse("Invalid Email or Password!"))
		c.Abort()
		return
	}

	td, err := services.CreateAuth(c, user.ID.Hex(), user.Role)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		c.Abort()
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
		c.Abort()
		return
	}

	deleted, err := services.DeleteAuth(c, accessUUID)
	if err != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Unauthorized!!"))
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("Successfully logged out", nil))
}

func RefreshToken(c *gin.Context) {
	var json map[string]string

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
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
		c.Abort()
		return
	}

	td, err := services.CreateAuth(c, claims.UserID, claims.Role)
	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		c.Abort()
		return
	}

	tokens := gin.H{
		"accessToken":  td.AccessToken,
		"refreshToken": td.RefreshToken,
	}
	c.JSON(http.StatusCreated, helpers.HandleSuccessResponse("", tokens))
}

func ChangePassword(c *gin.Context) {
	var json map[string]string

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	currentPass := json["currentPassword"]
	newPass := json["newPassword"]

	if currentPass == "" || newPass == "" {
		c.JSON(http.StatusUnprocessableEntity, helpers.HandleErrorResponse("`currentPassword` & `newPassword` are required!"))
		c.Abort()
		return
	}

	userID := c.GetString(helpers.CtxValues.UserID)
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, helpers.HandleErrorResponse(helpers.InvalidID))
		c.Abort()
		return
	}

	user, err := services.DBConn.FindUser(c, bson.M{"_id": objID, "isActive": true})
	if err != nil {
		c.JSON(http.StatusNotFound, helpers.HandleErrorResponse("Invalid Token!"))
		c.Abort()
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPass)); err != nil {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Invalid Password!"))
		c.Abort()
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(newPass), 14)
	if err != nil {
		fmt.Println("GenerateFromPassword err ==>", err.Error())
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	filters := bson.M{"_id": user.ID}

	update := bson.M{
		"$set": bson.M{
			"password": string(password),
		},
	}

	_, err = services.DBConn.UpdateUser(c, filters, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.HandleErrorResponse(err.Error()))
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, helpers.HandleSuccessResponse("Password Updated!", nil))
}
