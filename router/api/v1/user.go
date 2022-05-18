package v1

import (
	"net/http"
	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"
	"qiu/blog/pkg/util"
	service "qiu/blog/service"

	"github.com/gin-gonic/gin"
)

// Param
// 参数，表示需要传递到服务器端的参数，有五列参数，使用空格或者 tab 分割，五个分别表示的含义如下：
// 1.参数名
// 2.参数类型，可以有的值是 formData、query、path、body、header，formData 表示是 post 请求的数据，query 表示带在 url 之后的参数，path 表示请求路径上得参数，例如上面例子里面的 key，body 表示是一个 raw 数据请求，header 表示带在 header 信息中得参数。
// 3.参数类型
// 4.是否必须
// 5.注释

//用户登录
// @Summary 用户登录
// @Produce  json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {object} gin_http.ResponseJSON
// @Failure  20005 {object} gin_http.ResponseJSON
// @Router /login [post]
// @Content-type application/x-www-form-urlencoded
func Login(c *gin.Context) {
	userService := service.UserService{}
	httpCode, errCode := gin_http.Bind(c, &userService)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	// fmt.Println("绑定后", userService)
	state, err := userService.Login()
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_LOGIN, nil)
		return
	}
	if !state {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_LOGIN, nil)
		return
	}
	token, err := util.GenerateToken(userService.Username, userService.Password)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	data := make(map[string]interface{})
	data["token"] = token
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)

}

//用户注册
// @Summary 用户注册
// @Produce  json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {object} gin_http.ResponseJSON
// @Failure  20006 {object} gin_http.ResponseJSON
// @Failure  20007 {object} gin_http.ResponseJSON
// @Router /register [post]
func Register(c *gin.Context) {
	userService := service.UserService{}
	httpCode, errCode := gin_http.Bind(c, &userService)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	state := userService.IfExisted()

	if state {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_EXIST_USER, nil)
		return
	}

	err := userService.Register()
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_REGISTER, nil)
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}
