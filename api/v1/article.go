package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"qiu/blog/pkg/e"
	"qiu/blog/pkg/setting"
	"qiu/blog/pkg/upload"
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
// @Param pageNum query int false "Page"
// @Router /api/v1/article/list [get]
func GetArticles(c *gin.Context) {

	articleService := service.GetArticleService()
	page := 0
	articleService.PageNum, page = util.GetPage(c)
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
	data["datalist"] = articles
	data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = articleService.PageSize
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
	form, err := c.MultipartForm()
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_IMAGE_FAIL, nil)
		return
	}
	// 获取所有图片
	files := form.File["images"]
	// file, err := c.FormFile("image")
	// if err != nil {
	// 	fmt.Println(err)
	// 	gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_IMAGE_FAIL, nil)
	// 	return
	// }
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

		_, err = upload.Thumbnailify(imageName)
		if err != nil {
			fmt.Println(err.Error())
		}
		articleService.ImgName = append(articleService.ImgName, imageName)
	}

	// if err = upload.CheckImage(savePath); err != nil {
	// 	fmt.Println(err)
	// 	gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
	// 	return
	// }

	// httpCode, errCode = articleService.CheckTagName()
	// if errCode != e.SUCCESS {
	// 	gin_http.Response(c, httpCode, errCode, nil)
	// 	return
	// }

	// err := articleService.Add()
	err = articleService.AddArticleWithImg()
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
	if userID != claims.Uid && claims.Uid != setting.AppSetting.AdminId {
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

// @Summary 更新文章
// @Produce  json
// @Param id path int true "文章ID"
// @Param token header string true "token"
// @Router /api/v1/article/update/{id} [put]
func UpdateArticle(c *gin.Context) {

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
	if userID != claims.Uid && claims.Uid != setting.AppSetting.AdminId {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
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

// @Summary 添加文章点赞
// @Produce  json
// @Param id path uint true "文章ID"
// @Param user_id formData uint true "用户ID"
// @Param token header string true "token"
// @Router /api/v1/article/like/{id} [post]
func LikeArticle(c *gin.Context) {

	articleService := service.GetArticleService()
	param := service.ArticleLikeParams{}
	if err := c.ShouldBind(&param); err != nil {
		fmt.Println("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	param.Id, _ = strconv.Atoi(c.Param("id"))
	fmt.Println("绑定数据", param)
	claims := articleService.GetClaimsFromToken(c)
	if claims == nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}
	if param.UserID != claims.Uid {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH, nil)
		return
	}

	if err := articleService.AddArticleLikeUser(param); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_LIKE_ARTICLE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}
