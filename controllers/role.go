package controllers

import (
	"fmt"
	"net/http"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/gin-gonic/gin"
)

// func CreateRole(c *gin.Context) {
// 	var roleData []services.UserRole
// 	if err := c.ShouldBindJSON(&roleData); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	fmt.Println("roleData===>", roleData)
// 	data := make([]interface{}, len(roleData))
// 	for i, v := range roleData {
// 		data[i] = v
// 	}
// 	services.DBConn.InsertUserRole(data)
// }

func CreateRole(c *gin.Context) {
	var roleData services.UserRole
	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := services.DBConn.InsertUserRole(roleData); err != nil {
		fmt.Println("err ==>", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": helpers.CreatedMessage})
}
