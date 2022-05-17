package v1

import (
	"fmt"
	"net/http"
	"qiu/blog/model"
	"qiu/blog/pkg/e"
	"qiu/blog/pkg/logging"
	"qiu/blog/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func valid_username_password(username string, password string) bool {

	valid := validation.Validation{}

	valid.Required(username, "username").Message("用户名不能为空")

	valid.AlphaNumeric(username, "username").Message("用户名必须是数字或英文")

	valid.MaxSize(username, 10, "username").Message("用户名最长为10字符")

	valid.Required(password, "password").Message("密码不能为空")

	valid.MaxSize(password, 20, "password").Message("密码最长为20字符")

	valid.AlphaNumeric(password, "password").Message("密码必须是数字或英文")

	if valid.HasErrors() {
		for _, err := range valid.Errors {
			errf := fmt.Errorf("%q,+++++++,%q", err.Key, err.Message)
			logging.Info(errf)
			//log.Println(errf)
			//log.Println(err.Message)
		}
		return false

	}
	return true
}

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
// @Success 200 {string} string "ok"
// @Failure  400 {string} string "{"code":400,"data":{},"msg":"请求参数错误"}"
// @Failure  10000 {string} string "{"code":10001,"data":{},"msg":"TOKEN为空"}"
// @Failure  20003 {string} string "{"code":20003,"data":{},"msg":"Token生成失败"}"
// @Failure  20005 {string} string "{"code":20005,"data":{},"msg":"登录失败"}"
// @Router /login [post]
func Login(c *gin.Context) {
	content_type := c.ContentType()
	// fmt.Println("content-type", c.ContentType())
	var username = ""
	var password = ""
	if content_type == "application/x-www-form-urlencoded" {
		username = c.PostForm("username")
		password = c.PostForm("password")
	} else {
		json := make(map[string]string)
		c.BindJSON(&json)
		// fmt.Printf("%v", &json)
		username = json["username"]
		password = json["password"]
	}

	// fmt.Println("账户登录", username, password)
	data := make(map[string]interface{})
	code := e.INVALID_PARAMS

	if model.ValidLogin(username, password) && valid_username_password(username, password) {
		token, err := util.GenerateToken(username, password)
		if err != nil {
			code = e.ERROR_AUTH_TOKEN
		} else {
			data["token"] = token

			code = e.SUCCESS
		}
	} else {
		code = e.ERROR_LOGIN
	}
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//用户注册
// @Summary 用户注册
// @Produce  json
// @Param username query string true "username"
// @Param password query string true "password"
// @Success 200 {string} string "ok"
// @Failure  400 {string} string "{"code":400,"data":{},"msg":"输入数据有误"}"
// @Failure  10001 {string} string "{"code":10001,"data":{},"msg":"用户名已存在"}"
// @Router /register [post]
func AddUser(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if valid_username_password(username, password) {
		if model.ExistUsername(username) {
			c.JSON(http.StatusOK, gin.H{
				"code": e.ERROR_EXIST_TAG,
				"msg":  "用户名已存在",
				"data": make(map[string]string),
			})
		} else {
			model.AddUser(username, password, 1)
			c.JSON(http.StatusOK, gin.H{
				"code": e.SUCCESS,
				"msg":  "注册成功",
				"data": make(map[string]string),
			})
		}

	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": e.INVALID_PARAMS,
			"msg":  "输入数据有误",
			"data": make(map[string]string),
		})

	}

}
