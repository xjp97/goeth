package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func TestHandler() gin.HandlerFunc {

	return func(c *gin.Context) {
		fmt.Println("test handler")
		ass := c.Query("ass")
		if ass == "1" {
			fmt.Println("继续执行")

			// 继续执行
			c.Next()
		} else {
			// 中断
			fmt.Println("中断")
			c.Abort()
		}
		fmt.Println("来了老弟")
	}
}
