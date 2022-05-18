package gin_http

import (
	"qiu/blog/pkg/e"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpCode, errCode int, data interface{}) {
	c.JSON(httpCode, gin.H{
		"code": errCode,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})

}
