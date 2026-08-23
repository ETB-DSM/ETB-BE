package domain

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode,omitempty"`
}

func OK(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func Fail(c *gin.Context, err *AppError) {
	c.JSON(err.Status, ErrorResponse{
		Message:   err.Message,
		ErrorCode: err.ErrorCode,
	})
}
