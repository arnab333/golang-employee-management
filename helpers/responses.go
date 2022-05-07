package helpers

import "github.com/gin-gonic/gin"

func HandleErrorResponse(msg string) interface{} {
	var response = gin.H{
		"status": "error",
	}

	if msg != "" {
		response["message"] = msg
	}

	return response
}

func HandleSuccessResponse(msg string, data interface{}) interface{} {
	var response = gin.H{
		"status": "ok",
	}

	if data != nil {
		response["data"] = data
	}

	if msg != "" {
		response["message"] = msg
	}

	return response
}
