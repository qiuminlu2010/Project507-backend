package util

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/setting"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Uid uint `json:"uid"`
	// Username string `json:"username"`
	// Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(uid uint) (string, int64, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour).Unix()

	claims := Claims{
		uid,
		// username,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "qiu",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	fmt.Println("新建token缓存信息", token, claims)
	redis.Set(token, claims, 3600)
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
