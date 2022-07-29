package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"qiu/backend/pkg/e"
	"qiu/backend/pkg/redis"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		var token string
		code = e.SUCCESS

		//TODO: 重复TOKEN问题
		token = c.GetHeader("token")
		if token == "" {
			code = e.ERROR_AUTH
		} else {
			if !redis.SExist(e.CACHE_TOKEN, token) {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			}
			// claims, err := util.ParseToken(token)
			// if err != nil {
			// 	code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			// } else if time.Now().Unix() > claims.ExpiresAt {
			// 	code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			// }
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
