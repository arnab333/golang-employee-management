package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func VerifyToken(c *gin.Context) {
	bearToken := c.Request.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	var tokenString string
	if len(strArr) == 2 {
		tokenString = strArr[1]
	}

	claims, err := services.ExtractFromToken(tokenString, helpers.EnvKeys.JWT_ACCESS_SECRET)
	if err != nil {
		msg := "Invalid Token!"
		if strings.Contains(err.Error(), "expired") {
			msg = "Token Expired!"
		}
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse(msg))
		c.Abort()
		return
	}

	if claims.ID == "" {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Invalid Token!!"))
		c.Abort()
		return
	}

	userID, err := services.FetchAuth(c, claims.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Invalid Token!!!"))
		c.Abort()
		return
	}

	c.Set(helpers.CtxValues.UserID, userID)
	c.Set(helpers.CtxValues.AccessUUID, claims.ID)
	c.Set(helpers.CtxValues.Role, claims.Role)
	c.Next()
}

func CheckRole(c *gin.Context) {
	roleCtx := c.GetString(helpers.CtxValues.Role)
	objID, err := primitive.ObjectIDFromHex(roleCtx)

	if err != nil {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("You are not Authorized!"))
		c.Abort()
		return
	}

	filters := bson.M{"_id": objID}
	role, err := services.DBConn.FindRole(c, filters)
	if err != nil {
		fmt.Println("FindRole err ===>", err.Error())
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("You are not Authorized!!"))
		c.Abort()
		return
	}

	c.Set(helpers.CtxValues.Permissions, role.Permissions)
	c.Next()
}
