package hello

import "github.com/gin-gonic/gin"

// Hello 测试接口
// @Summary 测试接口
// @Description 测试接口
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /example/hello_world [get]
func Hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message_type": "Hello, World!",
	})
}
