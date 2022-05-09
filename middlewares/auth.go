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
	c.Next()
}