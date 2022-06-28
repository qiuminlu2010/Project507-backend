package v1

import (
	"fmt"
	"net/http"
	"qiu/blog/pkg/e"
	gin_http "qiu/blog/pkg/http"
	service "qiu/blog/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
	if err := c.ShouldBind(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	params.UserId, _ = strconv.Atoi(c.Param("id"))
	fmt.Println("绑定数据", params)
	if !userService.CheckTokenUid(c, uint(params.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if err := userService.UpsertFollowUser(params); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.EEROR_USER_UPSERT_FOLLOW_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 关注列表
// @Produce  json
// @Param id path uint true "用户ID"
// @Param token header string true "token"
// @Router /api/v1/user/{id}/follow [get]
func GetFollows(c *gin.Context) {
	userService := service.GetUserService()
	params := service.UserFollowsParams{}
	params.UserId, _ = strconv.Atoi(c.Param("id"))
	fmt.Println("绑定数据", params)
	if !userService.CheckTokenUid(c, uint(params.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	data, err := userService.GetFollows(params)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.EEROR_USER_UPSERT_FOLLOW_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 用户动态列表
// @Produce  json
// @Param id path uint true "用户ID"
// @Param token header string true "token"
// @Router /api/v1/user/{id}/follow [get]
func GetUserArticles(c *gin.Context) {
	userService := service.GetUserService()
	params := service.UserFollowsParams{}
	params.UserId, _ = strconv.Atoi(c.Param("id"))
	fmt.Println("绑定数据", params)
	if !userService.CheckTokenUid(c, uint(params.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	data, err := userService.GetFollows(params)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.EEROR_USER_UPSERT_FOLLOW_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 用户喜欢列表
// @Produce  json
// @Param id path uint true "用户ID"
// @Param token header string true "token"
// @Router /api/v1/user/{id}/follow [get]
func GetUserLikeArticles(c *gin.Context) {
	userService := service.GetUserService()
	params := service.UserFollowsParams{}
	params.UserId, _ = strconv.Atoi(c.Param("id"))
	fmt.Println("绑定数据", params)
	if !userService.CheckTokenUid(c, uint(params.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	data, err := userService.GetFollows(params)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.EEROR_USER_UPSERT_FOLLOW_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}
