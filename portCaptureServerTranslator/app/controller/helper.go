package controller

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func getSuccessResponse() *response {
	return &response{
		Success: true,
	}
}

func getFailResponse(err string) *response {
	return &response{
		Success: false,
		Error:   err,
	}
}

func abortAndError(c *gin.Context, err error) {
	c.Abort()

	c.JSON(400, getFailResponse(err.Error()))
}
