package middlewares

import (
	"net/http"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func VerifyRole(c *gin.Context) {
	roleCtx := c.GetString(helpers.CtxValues.Role)

	filters := bson.M{"name": roleCtx}
	role, err := services.DBConn.FindRole(c, filters)

	if err != nil {
		c.JSON(http.StatusUnauthorized, helpers.HandleErrorResponse("You are not Authorized!"))
		c.Abort()
		return
	}

	c.Set(helpers.CtxValues.Permissions, role.Permissions)
	c.Next()
}
