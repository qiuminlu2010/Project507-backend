package v1

import (
	"fmt"
	"net/http"
	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"
	"qiu/blog/pkg/setting"
	service "qiu/blog/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 搜索用户
// @Produce  json
// @Param name query string true "用户名"
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Router /api/v1/search/user [get]
func GetUsers(c *gin.Context) {

	userService := service.GetUserService()
	params := service.UsersGetParams{}

	if err := c.ShouldBind(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}
	page := params.PageNum
	params.PageNum = params.PageNum * params.PageSize

	fmt.Println("绑定数据", params)

	users, err := userService.GetUsersByName(&params)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_USER_LIST_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["datalist"] = users
	// data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = params.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)

}

// @Summary 获取用户信息
// @Produce  json
// @Param id path uint true "用户ID"
// @Router /api/v1/user/{id} [get]
func GetUserInfo(c *gin.Context) {
	userService := service.GetUserService()
	var err error
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	data, err := userService.GetUserInfo(userId)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_USER_INFO_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 关注用户
// @Produce  json
// @Param id path uint true "用户ID"
// @Param follow_id formData int true "关注用户ID"
// @Param type formData int true "类型"
// @Param token header string true "token"
// @Router /api/v1/user/{id}/follow [post]
func FollowUser(c *gin.Context) {
	userService := service.GetUserService()
	params := service.UpsertUserFollowParams{}

	var err error
	if err = c.ShouldBind(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	params.UserId, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	fmt.Println("绑定数据", params)
	if !userService.CheckTokenUid(c, uint(params.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if err := userService.UpsertFollowUser(params); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_USER_UPSERT_FOLLOW_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 关注列表
// @Produce  json
// @Param id path uint true "用户ID"
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Param token header string true "token"
// @Router /api/v1/user/{id}/follows [get]
func GetFollows(c *gin.Context) {
	userService := service.GetUserService()
	params := service.FollowsGetParams{}

	var err error
	params.UserId, err = strconv.Atoi(c.Param("id"))

	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if err := c.ShouldBind(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}

	page := params.PageNum
	params.PageNum = params.PageNum * params.PageSize

	fmt.Println("绑定数据", params)

	if !userService.CheckTokenUid(c, uint(params.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	follows, err := userService.GetFollows(&params)

	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_USER_UPSERT_FOLLOW_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["datalist"] = follows
	// data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = params.PageSize

	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 粉丝列表
// @Produce  json
// @Param id path uint true "用户ID"
// @Param token header string true "token"
// @Router /api/v1/user/{id}/fans [get]
func GetFans(c *gin.Context) {
	userService := service.GetUserService()

	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !userService.CheckTokenUid(c, uint(userId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	data, err := userService.GetFans(userId)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_USER_UPSERT_FOLLOW_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 用户动态列表
// @Produce  json
// @Param id path uint true "用户ID"
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Router /api/v1/user/{id}/articles [get]
func GetUserArticles(c *gin.Context) {

	userService := service.GetUserService()
	params := service.ArticleGetParams{}

	var err error
	params.Uid, err = strconv.Atoi(c.Param("id"))

	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}

	page := params.PageNum
	params.PageNum = params.PageNum * params.PageSize

	fmt.Println("绑定数据", params)

	articles, err := userService.GetUserArticles(&params)
	if err != nil {
		fmt.Println("GetUserLikeArticles", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_USER_ARTICLES_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["data"] = articles
	// data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = params.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 用户喜欢列表
// @Produce  json
// @Param id path uint true "用户ID"
// @Param token header string true "token"
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Router /api/v1/user/{id}/likeArticles [get]
func GetUserLikeArticles(c *gin.Context) {
	userService := service.GetUserService()
	params := service.ArticleGetParams{}
	params.Uid, _ = strconv.Atoi(c.Param("id"))

	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}
	page := params.PageNum
	params.PageNum = params.PageNum * params.PageSize
	fmt.Println("绑定数据", params)
	if !userService.CheckTokenUid(c, uint(params.Uid)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	articles, err := userService.GetLikeArticles(&params)
	if err != nil {
		fmt.Println("GetUserLikeArticles", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_LIKE_ARTICLES_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["data"] = articles
	// data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = params.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}
