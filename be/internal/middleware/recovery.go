package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery Gin中间件，恢复panic并记录堆栈跟踪
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				fmt.Printf("[RECOVERY] %v recovered\n%s\n", err, stack)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"message": "Internal Server Error",
					"data":    nil,
				})
			}
		}()
		c.Next()
	}
}
