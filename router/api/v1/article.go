package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/e"
	"qiu/blog/pkg/setting"
	"qiu/blog/pkg/util"

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
// @Param page query int false "Page"
// @Router /api/v1/article/getList [get]
func GetArticles(c *gin.Context) {

	articleService := service.GetArticleService()
	articleService.PageNum = util.GetPage(c)
	articleService.PageSize = setting.AppSetting.PageSize

	// data["delete_on"] = 0
	// data["state"] = 1
	total, err := articleService.Count(nil)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetArticles(nil)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total

	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 添加文章
// @Produce  json
// @Param user_id formData int true "用户id"
// @Param content formData string true "内容"
// @Param tag_name formData []string false "标签"
// @Param token header string true "token"
// @Router /api/v1/article/add [post]
func AddArticle(c *gin.Context) {

	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := articleService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	//TODO:还需验证用户是否存在
	if articleService.UserID != claims.Uid {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	httpCode, errCode = articleService.CheckTagName()
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	err := articleService.Add()
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}

// @Summary 添加文章标签
// @Produce  json
// @Param id path int true "文章ID"
// @Param tag_id formData []uint false "标签ID"
// @Param token header string true "token"
// @Router /api/v1/article/addTags/{id} [post]
func AddArticleTags(c *gin.Context) {

	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := articleService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}
	userID, err := articleService.GetUserID()
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if userID != claims.Uid {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
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
// @Router /api/v1/article/deleteTags/{id} [delete]
func DeleteArticleTags(c *gin.Context) {

	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := articleService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}
	userID, err := articleService.GetUserID()
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if userID != claims.Uid && claims.Uid != 1 {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
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
// @Router /api/v1/article/delete/{id} [delete]
func DeleteArticle(c *gin.Context) {

	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := articleService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}
	userID, err := articleService.GetUserID()
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if userID != claims.Uid && claims.Uid != 1 {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	if err := articleService.Delete(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 恢复文章
// @Produce  json
// @Param id path int true "文章ID"
// @Param token header string true "token"
// @Router /api/v1/article/recover/{id} [post]
func RecoverArticle(c *gin.Context) {
	articleService := service.GetArticleService()
	httpCode, errCode := articleService.Bind(c)
	if errCode != e.SUCCESS {
		gin_http.Response(c, httpCode, errCode, nil)
		return
	}

	claims := articleService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}
	userID, err := articleService.GetUserID()
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if userID != claims.Uid {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	if err := articleService.Recovery(); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_REC_ARTICLE_FAIL, nil)
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}
