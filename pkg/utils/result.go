package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseStatusCode int

const (
	SuccessCode ResponseStatusCode = iota
	ErrorCode
	FailCode
)

type Response struct {
	Status  ResponseStatusCode `json:"status"`         // 状态码
	Message string             `json:"message"`        // 提示信息
	Data    interface{}        `json:"data,omitempty"` // 返回数据（可以为空）
}

// Success 成功返回
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Status:  SuccessCode,
		Message: message,
		Data:    data,
	})
	c.Abort()
}

// SuccessWithDefault 成功返回，默认信息
func SuccessWithDefault(c *gin.Context, data interface{}) {
	Success(c, data, "success")
	c.Abort()
}

// Error 错误返回
func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Status:  ErrorCode,
		Message: message,
	})
	c.Abort()
}

// ErrorWithDefault 错误返回，默认信息
func ErrorWithDefault(c *gin.Context) {
	Error(c, "error")
	c.Abort()
}

// Fail 失败返回
func Fail(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Status:  FailCode,
		Message: message,
		Data:    data,
	})
	c.Abort()
}

// FailWithDefault 失败返回，默认信息
func FailWithDefault(c *gin.Context, data interface{}) {
	Fail(c, data, "fail")
	c.Abort()
}
