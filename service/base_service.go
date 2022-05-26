package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"qiu/blog/pkg/e"
	"qiu/blog/pkg/redis"
	"qiu/blog/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type BaseService struct {
	model interface{}
}

func (s *BaseService) Bind(c *gin.Context) (int, int) {
	var err error
	// fmt.Println("绑定json")
	// if err = c.ShouldBindJSON(s.model); err != nil {
	// 	// fmt.Println("绑定数据", s.model)
	// 	// fmt.Println("绑定错误", err)
	// 	return http.StatusBadRequest, e.INVALID_PARAMS
	// }
	fmt.Println("绑定参数")
	if err = c.ShouldBind(s.model); err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	fmt.Println("绑定url")
	if err = c.ShouldBindUri(s.model); err != nil {
		// fmt.Println("绑定数据", s.model)
		// fmt.Println("绑定错误", err)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}
	fmt.Println("绑定数据", s.model)
	return http.StatusOK, e.SUCCESS
}

func (s *BaseService) Valid() error {
	validate := validator.New()
	err := validate.Struct(s.model)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("%v should %v %v, but got %v\n", err.Namespace(), err.Tag(), err.Param(), err.Value())
			// logging.Info("%v should %v %v, but got %v", err.Namespace(), err.Tag(), err.Param(), err.Value())
		}
		return err
	}
	return nil
}

func (s *BaseService) GetClaimsFromToken(c *gin.Context) *util.Claims {

	claims := &util.Claims{}
	token := c.GetHeader("token")

	if token == "" {
		return nil
	}
	if redis.Exists(token) {
		data, err := redis.Get(token)
		if err != nil {
			return nil
		} else {
			json.Unmarshal(data, &claims)
		}
		fmt.Println("获取token缓存信息", claims)
		return claims
	}
	claims, err := util.ParseToken(token)
	if err != nil {
		return nil
	}
	fmt.Println("新建token缓存信息", token, claims)
	err = redis.Set(token, claims, 3600)
	if err != nil {
		fmt.Println(err)
	}
	return claims

}
