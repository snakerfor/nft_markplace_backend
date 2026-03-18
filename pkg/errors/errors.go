package errors

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.Code, appErr.Message)
		return
	}

	// 未知错误
	Error(c, 500, "Internal server error")
}

// Error 返回错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"error":   message,
	})
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    data,
	})
}

// ValidationError 返回验证错误响应
func ValidationError(c *gin.Context, errs map[string]string) {
	c.JSON(422, gin.H{
		"code":    422,
		"message": "validation failed",
		"error":   errs,
	})
}
