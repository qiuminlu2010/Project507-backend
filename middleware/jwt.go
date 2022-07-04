package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/e"
	"qiu/blog/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		var token string
		code = e.SUCCESS

		token = c.GetHeader("token")
		if token == "" {
			code = e.ERROR_AUTH
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
			// } else if redis.Exists(token) == 0 {
			// 	fmt.Println("新建Redis缓存", token, claims)
			// 	cacha_key := service.GetKeyName("user", claims.Uid, "token")
			// 	redis.Set(cacha_key, token, claims.TTL)
			// }
		}

		if code != e.SUCCESS {
			// logging.Info(e.GetMsg(code))
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
