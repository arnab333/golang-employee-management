package middlewares

import (
	"net/http"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
)

func VerifyToken(c *gin.Context) {
	bearToken := c.Request.Header.Get("Authorization")
	claims, err := services.ExtractFromToken(bearToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Invalid Token!"))
		c.Abort()
		return
	}

	// services.AuthDetails.UserID = claims["iss"].(string)

	accessUUID, ok := claims["jti"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Invalid Token!!"))
		c.Abort()
		return
	}

	userID, err := services.FetchAuth(accessUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("Invalid Token!!!"))
		c.Abort()
		return
	}

	c.Set(helpers.CtxValues.UserID, userID)
	c.Set(helpers.CtxValues.AccessUUID, accessUUID)
	c.Next()
}
