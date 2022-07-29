package v1

import (
	"net/http"
	"qiu/backend/pkg/e"
	gin_http "qiu/backend/pkg/http"
	log "qiu/backend/pkg/logging"
	"qiu/backend/pkg/setting"
	"qiu/backend/pkg/upload"
	articleService "qiu/backend/service/article"
	param "qiu/backend/service/param"
	service "qiu/backend/service/user"
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
	params := param.UsersGetParams{}

	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}
	page := params.PageNum
	params.PageNum = params.PageNum * params.PageSize
	log.Logger.Debug("绑定数据", params)

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
	params := param.UpsertUserFollowParams{}

	var err error
	if err = c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	params.UserId, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	log.Logger.Debug("绑定数据", params)
	if !userService.CheckTokenUid(c, uint(params.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if err := userService.UpsertFollowUser(&params); err != nil {
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
	params := param.FollowsGetParams{}

	var err error
	params.UserId, err = strconv.Atoi(c.Param("id"))

	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}

	page := params.PageNum
	params.PageNum = params.PageNum * params.PageSize

	log.Logger.Debug("绑定数据", params)

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

	articleService := articleService.GetArticleService()
	params := param.ArticleGetParams{}

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

	log.Logger.Debug("绑定数据", params)

	articles, err := articleService.GetUserArticles(&params)
	if err != nil {
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
	articleService := articleService.GetArticleService()
	params := param.ArticleGetParams{}
	params.Uid, _ = strconv.Atoi(c.Param("id"))

	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}
	page := params.PageNum
	params.PageNum = params.PageNum * params.PageSize
	log.Logger.Debug("绑定数据", params)
	if !articleService.CheckTokenUid(c, uint(params.Uid)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	articles, err := articleService.GetUserLikeArticles(&params)
	if err != nil {
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

// @Summary 更新用户头像
// @Produce  json
// @Param id path uint true "用户ID"
// @Param image formData file true "image"
// @Param token header string true "token"
// @Router /api/v1/user/{id}/avatar [put]
func UpdateAvatar(c *gin.Context) {
	userService := service.GetUserService()
	var err error
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil || userId <= 0 {
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if !userService.CheckTokenUid(c, uint(userId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	_, image, err := c.Request.FormFile("img")
	if err != nil {
		log.Logger.Error("保存图片失败", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	if image == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_IMAGE_LOST, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	savePath := upload.GetAvatarSavePath()

	src := savePath + imageName
	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(image) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}
	log.Logger.Info("保存上传图片", src)
	if err = c.SaveUploadedFile(image, "."+src); err != nil {
		log.Logger.Error("保存图片失败", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}
	if err := userService.UpdateAvatar(userId, src); err != nil {
		log.Logger.Error("更新头像失败", err)
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPDATE_USER_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["avatar"] = src
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}
