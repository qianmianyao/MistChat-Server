package hello

import "github.com/gin-gonic/gin"

// Hello handles the hello world request
// @Summary Hello World
// @Description 返回 Hello World
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} string "Hello World"
// @Router /example/hello_world [get]
func Hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message_type": "Hello, World!",
	})
}
