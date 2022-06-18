package v1

import (
	"fmt"
	"net/http"
	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"
	"qiu/blog/pkg/logging"
	"qiu/blog/pkg/util"
	service "qiu/blog/service"

	"qiu/blog/pkg/setting"

	"github.com/gin-gonic/gin"
)

// Param
// 参数，表示需要传递到服务器端的参数，有五列参数，使用空格或者 tab 分割，五个分别表示的含义如下：
// 1.参数名
// 2.参数类型，可以有的值是 formData、query、path、body、header，formData 表示是 post 请求的数据，query 表示带在 url 之后的参数，path 表示请求路径上得参数，例如上面例子里面的 key，body 表示是一个 raw 数据请求，header 表示带在 header 信息中得参数。
// 3.参数类型
// 4.是否必须
// 5.注释

// @Summary 用户列表
// @Produce  json
// @Param token header string true "token"
// @Success 200 {object} gin_http.ResponseJSON
// @Router /user/list [get]
func GetUserList(c *gin.Context) {

	userService := service.GetUserService()
	claims := userService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	if claims.Uid != userService.Id && claims.Uid != setting.AppSetting.AdminId {
		fmt.Println("token用户信息不一致", claims.Uid, userService.Id)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	userService.PageNum = util.GetPage(c)
	userService.PageSize = setting.AppSetting.PageSize

	total, err := userService.CountUser(nil)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_USER_LIST_FAIL, nil)
		return
	}
	userList, err := userService.GetUserList(nil)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_USER_LIST_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["datalist"] = userList
	data["total"] = total
	data["pageNum"] = userService.PageNum
	data["pageSize"] = userService.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 用户登录
// @Produce  json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Success 200 {object} gin_http.ResponseJSON
// @Failure  20005 {object} gin_http.ResponseJSON
// @Router /user/login [post]
func Login(c *gin.Context) {
	userService := service.GetUserService()
	httpCode, errCode := userService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	userInfo, err := userService.Login()
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_LOGIN, nil)
		return
	}

	token, expire_time, err := util.GenerateToken(userInfo.ID)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	data := make(map[string]interface{})
	data["uid"] = userInfo.ID
	data["username"] = userInfo.Username
	data["token"] = token
	data["uuid"] = userService.GetUUID(userInfo.ID)
	data["expire_time"] = expire_time
	logging.Info("用户登录成功,", "用户名:", userService.Username)
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
// @Router /user/register [post]
func Register(c *gin.Context) {

	userService := service.GetUserService()
	httpCode, errCode := userService.Bind(c)

	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	if err := userService.ExistUsername(); err == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_EXIST_USER, nil)
		return
	}
	//TODO:密码加密
	if err := userService.Add(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_REGISTER, nil)
		return
	}

	logging.Info("用户注销成功,", "用户名:", userService.Username)
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 注销用户
// @Produce  json
// @Param id path int true "id"
// @Param username formData string true "username"
// @Param token header string true "token"
// @Success 200 {object} gin_http.ResponseJSON
// @Failure  400 {object} gin_http.ResponseJSON
// @Failure  20008 {object} gin_http.ResponseJSON
// @Router /user/delete/{id} [delete]
func DeleteUser(c *gin.Context) {

	userService := service.GetUserService()
	httpCode, errCode := userService.Bind(c)

	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := userService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	if claims.Uid != userService.Id && claims.Uid != setting.AppSetting.AdminId {
		fmt.Println("token用户信息不一致", claims.Uid, userService.Id)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := userService.Delete(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_DELETE_USER_FAIL, nil)
		return
	}
	logging.Info("用户注销成功,", "用户名:", userService.GetUsernameByID())
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 修改用户密码
// @Produce  json
// @Param id path int true "id"
// @Param password formData string true "password"
// @Param token header string true "token"
// @Success 200 {object} gin_http.ResponseJSON
// @Failure  400 {object} gin_http.ResponseJSON
// @Failure  20009 {object} gin_http.ResponseJSON
// @Router /user/update/{id} [put]
func UpdatePassword(c *gin.Context) {

	userService := service.GetUserService()
	httpCode, errCode := userService.Bind(c)

	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := userService.GetClaimsFromToken(c)
	// token := c.GetHeader("token")
	// claims, err := util.ParseToken(token)

	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	if claims.Uid != userService.Id && claims.Uid != setting.AppSetting.AdminId {
		fmt.Println("token用户信息不一致", claims.Uid, userService.Id)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := userService.UpdatePassword(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPDATE_USER_FAIL, nil)
		return
	}
	logging.Info("用户修改密码成功,", "用户名:", userService.GetUsernameByID())
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 更新Token
// @Produce  json
// @Param id path uint true "id"
// @Param uuid formData string true "uuid"
// @Success 200 {object} gin_http.ResponseJSON
// @Router /user/refreshToken/{id} [post]
func RefreshToken(c *gin.Context) {
	userService := service.GetUserService()
	httpCode, errCode := userService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	uuid := c.Param("uuid")
	state := userService.CheckUUID(userService.Id, uuid)
	if !state {
		gin_http.Response(c, http.StatusMovedPermanently, e.ERROR_UUID_EXPIRE, nil)
		return
	}
	token, expire_time, err := util.GenerateToken(userService.Id)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}
	data := make(map[string]interface{})
	// data["uid"] = userService.Id
	data["token"] = token
	data["expire_time"] = expire_time
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)

}

// @Summary 更新用户
// @Produce  json
// @Param id path int true "id"
// @Param token header string true "token"
// @Success 200 {object} gin_http.ResponseJSON
// @Router /user/update/{id} [delete]
func UpdateUserState(c *gin.Context) {

	userService := service.GetUserService()
	httpCode, errCode := userService.Bind(c)

	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := userService.GetClaimsFromToken(c)
	// token := c.GetHeader("token")
	// claims, err := util.ParseToken(token)

	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	if claims.Uid != userService.Id && claims.Uid != setting.AppSetting.AdminId {
		fmt.Println("token用户信息不一致", claims.Uid, userService.Id)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := userService.UpdateState(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPDATE_USER_FAIL, nil)
		return
	}
	// logging.Info("用户修改密码成功,", "用户名:", userService.GetUsernameByID())
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}
