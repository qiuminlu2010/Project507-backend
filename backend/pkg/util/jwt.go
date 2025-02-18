package util

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"qiu/backend/pkg/e"
	"qiu/backend/pkg/redis"
	"qiu/backend/pkg/setting"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Uid uint          `json:"uid"`
	TTL time.Duration `json:"ttl"`
	// Username string `json:"username"`
	// Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(uid uint) (string, int64, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(e.DURATION_USER_TOKEN).Unix()
	ttl := e.DURATION_USER_TOKEN
	claims := Claims{
		uid,
		ttl,
		// username,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "qiu",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	fmt.Println("新建token缓存信息", token, claims)
	cacha_key := fmt.Sprintf("%s:%d:%s", "user", claims.Uid, "token")

	//TODO: 事务 lua脚本
	if redis.Exists(cacha_key) != 0 {
		old := redis.Get(cacha_key)
		redis.SDEL(e.CACHE_TOKEN, old)
	}
	redis.Set(cacha_key, token, claims.TTL)
	redis.SAdd(e.CACHE_TOKEN, token)
	// redis.Set(token, cache, 3600)
	return token, expireTime, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
