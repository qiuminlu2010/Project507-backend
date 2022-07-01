package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/e"
	"qiu/blog/pkg/setting"
	"qiu/blog/pkg/upload"

	gin_http "qiu/blog/pkg/http"

	service "qiu/blog/service"
)

//获取单个文章
func GetArticle(c *gin.Context) {

	// articleService := service.ArticleService{}
	// httpCode, errCode := gin_http.Bind(c, &articleService)

	// if errCode != e.SUCCESS {
	// 	gin_http.Response(c, httpCode, errCode, nil)
	// 	return
	// }

	// exists, err := articleService.ExistByID()

	// if err != nil {
	// 	gin_http.Response(c, http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
	// 	return
	// }
	// if !exists {
	// 	gin_http.Response(c, http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
	// 	return
	// }
	// getArticle, err := articleService.Get()
	// if err != nil {
	// 	gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
	// 	return
	// }

	// gin_http.Response(c, http.StatusOK, e.SUCCESS, getArticle)
}

// @Summary 获取文章列表
// @Produce  json
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Param uid query int false "uid"
// @Router /api/v1/article [get]
func GetArticles(c *gin.Context) {

	articleService := service.GetArticleService()
	params := service.ArticleGetParams{}
	// articleService.PageNum, page = util.GetPage(c)
	// articleService.PageSize = setting.AppSetting.PageSize
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

	articles, err := articleService.GetArticles(params)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	total := articleService.Count(nil)
	// if err != nil {
	// 	gin_http.Response(c, http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
	// 	return
	// }
	data := make(map[string]interface{})
	data["datalist"] = articles
	data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = params.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 添加文章
// @Produce  json
// @Accept multipart/form-data
// @Param user_id formData int true "用户id"
// @Param content formData string true "内容"
// @Param tag_name formData []string false "标签"
// @Param images formData file true "image"
// @Param token header string true "token"
// @Router /api/v1/article [post]
func AddArticle(c *gin.Context) {
	articleService := service.GetArticleService()
	params := service.ArticleAddParams{}
	if err := c.ShouldBind(&params); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if !articleService.CheckTokenUid(c, params.UserID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_IMAGE_FAIL, nil)
		return
	}

	// 获取所有图片
	files := form.File["images"]

	for _, file := range files {
		imageName := upload.GetImageName(file.Filename)
		savePath := upload.GetImagePath() + imageName
		fmt.Println("文件名", imageName)
		fmt.Println("保存路径", savePath)
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
			return
		}
		if err = c.SaveUploadedFile(file, savePath); err != nil {
			fmt.Println(err)
			gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
			return
		}

		params.ImgName = append(params.ImgName, imageName)
	}

	err = articleService.AddArticleWithImg(&params)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

	for _, imageName := range params.ImgName {
		//TODO: 异步
		go func(s string) {
			_, err = upload.Thumbnailify(s)
			if err != nil {
				fmt.Println(err.Error())
			}
		}(imageName)
	}

}

// @Summary 添加文章标签
// @Produce  json
// @Param id path int true "文章ID"
// @Param tag_id formData []uint false "标签ID"
// @Param token header string true "token"
// @Router /api/v1/article/{id}/addTags [post]
func AddArticleTags(c *gin.Context) {

	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	articleId, _ := strconv.Atoi(c.Param("id"))
	userID, err := articleService.GetUserID(articleId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if !articleService.CheckTokenUid(c, userID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if httpCode, errCode = articleService.AddArticleTags(); errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}

// @Summary 删除文章标签
// @Produce  json
// @Param id path int true "文章ID"
// @Param tag_id formData []int true "标签ID"
// @Param token header string true "token"
// @Router /api/v1/article/{id}/deleteTags [delete]
func DeleteArticleTags(c *gin.Context) {

	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	articleId, _ := strconv.Atoi(c.Param("id"))
	userID, err := articleService.GetUserID(articleId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if !articleService.CheckTokenUid(c, userID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if httpCode, errCode = articleService.DeleteArticleTags(); errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}

// @Summary 删除文章
// @Produce  json
// @Param id path int true "文章ID"
// @Param token header string true "token"
// @Router /api/v1/article/{id} [delete]
func DeleteArticle(c *gin.Context) {

	articleService := service.GetArticleService()
	// httpCode, errCode := articleService.Bind(c)
	// if errCode != e.SUCCESS {
	// 	gin_http.Response(c, httpCode, errCode, nil)
	// 	return
	// }
	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil || articleId <= 0 {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	fmt.Println("绑定数据", articleId)
	userID, err := articleService.GetUserID(articleId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if !articleService.CheckTokenUid(c, userID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := articleService.Delete(articleId); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 更新文章
// @Produce  json
// @Param id path int true "文章ID"
// @Param token header string true "token"
// @Router /api/v1/article/{id} [put]
func UpdateArticle(c *gin.Context) {

	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	articleId, _ := strconv.Atoi(c.Param("id"))
	userID, err := articleService.GetUserID(articleId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if !articleService.CheckTokenUid(c, userID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := articleService.Update(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 恢复文章
// @Produce  json
// @Param id path int true "文章ID"
// @Param token header string true "token"
// @Router /api/v1/article/{id}/recover [post]
func RecoverArticle(c *gin.Context) {
	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}
	articleId, _ := strconv.Atoi(c.Param("id"))
	userID, err := articleService.GetUserID(articleId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if !articleService.CheckTokenUid(c, userID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := articleService.Recovery(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_REC_ARTICLE_FAIL, nil)
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 添加文章点赞
// @Produce  json
// @Param id path uint true "文章ID"
// @Param user_id formData uint true "用户ID"
// @Param type formData int true "类型"
// @Param token header string true "token"
// @Router /api/v1/article/{id}/like [post]
func LikeArticle(c *gin.Context) {

	articleService := service.GetArticleService()
	param := service.ArticleLikeParams{}
	if err := c.ShouldBind(&param); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	var err error
	param.ArticleId, err = strconv.Atoi(c.Param("id"))
	if err != nil || param.ArticleId <= 0 {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	fmt.Println("绑定数据", param)
	if !articleService.CheckTokenUid(c, uint(param.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if err := articleService.UpdateArticleLike(param); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_LIKE_ARTICLE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}
