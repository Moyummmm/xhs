package response

import (
	"github.com/gin-gonic/gin"
)

// Success 成功响应（带数据）
func Success(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": data,
	})
}

// SuccessMsg 成功响应（带自定义消息）
func SuccessMsg(c *gin.Context, msg string, data any) {
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  msg,
		"data": data,
	})
}

// Fail 失败响应（业务错误码）
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
	})
}

// FailWithStatus 带 HTTP 状态码的失败响应
func FailWithStatus(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, gin.H{
		"code": code,
		"msg":  msg,
	})
}
