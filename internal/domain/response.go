package domain

import "github.com/gin-gonic/gin"

type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Message   string      `json:"message,omitempty"`
	ErrorCode string      `json:"errorCode,omitempty"`
}

func OK(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Response{Success: true, Data: data})
}

func Fail(c *gin.Context, err *AppError) {
	c.JSON(err.Status, Response{
		Success:   false,
		Message:   err.Message,
		ErrorCode: err.ErrorCode,
	})
}
