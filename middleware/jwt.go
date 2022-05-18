package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/e"
	"qiu/blog/pkg/logging"
	"qiu/blog/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.ERROR_AUTH
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			logging.Info(e.GetMsg(code))
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
