package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`         // 状态码
	Message string      `json:"message"`        // 提示信息
	Data    interface{} `json:"data,omitempty"` // 返回数据（可以为空）
}

// Success 成功返回
func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, Response{
		Status:  0,
		Message: message,
		Data:    data,
	})
}

// SuccessWithDefault 成功返回，默认信息
func SuccessWithDefault(c *gin.Context, data interface{}) {
	Success(c, data, "success")
}

// Error 错误返回
func Error(c *gin.Context, status int, message string) {
	c.JSON(http.StatusOK, Response{
		Status:  status,
		Message: message,
	})
}
