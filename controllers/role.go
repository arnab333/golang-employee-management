package controllers

import (
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
)

func CreateRole(ctx *gin.Context) {
	roleData := []services.UserRole{
		{
			Value: 1,
			Text:  "admin",
		},
	}

	data := make([]interface{}, len(roleData))
	for i, v := range roleData {
		data[i] = v
	}
	services.DBConn.InsertUserRole(data)
}
