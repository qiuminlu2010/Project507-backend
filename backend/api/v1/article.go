package v1

import (
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"

	"qiu/backend/pkg/e"
	gin_http "qiu/backend/pkg/http"
	log "qiu/backend/pkg/logging"
	"qiu/backend/pkg/minio"
	"qiu/backend/pkg/setting"
	"qiu/backend/pkg/upload"
	service "qiu/backend/service/article"
	param "qiu/backend/service/param"
)

// @Summary 获取文章列表
// @Produce  json
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Param uid query int false "uid"
// @Router /api/v1/article [get]
func GetArticles(c *gin.Context) {

	articleService := service.GetArticleService()
	params := param.ArticleGetParams{}
	// articleService.PageNum, page = util.GetPage(c)
	// articleService.PageSize = setting.AppSetting.PageSize
	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	if params.PageSize == 0 {
		params.PageSize = setting.AppSetting.PageSize
	}
	// page := params.PageNum
	// params.PageNum = params.PageNum * params.PageSize
	log.Logger.Debug("绑定数据", params)

	articles, err := articleService.GetArticles(&params)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	total := articleService.Count(nil)

	data := make(map[string]interface{})
	data["datalist"] = articles
	data["total"] = total
	data["pageNum"] = params.PageNum
	data["pageSize"] = params.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 添加文章
// @Produce  json
// @Accept multipart/form-data
// @Param user_id formData int true "用户id"
// @Param content formData string true "内容"
// @Param tag_name formData []int false "标签"
// @Param images formData file false "image"
// @Param video formData file false "video"
// @Param type formData int false "视频类型为1"
// @Param token header string true "token"
// @Router /api/v1/article [post]
func AddArticle(c *gin.Context) {
	articleService := service.GetArticleService()
	params := param.ArticleAddParams{}
	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
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

	// 视频类型
	if params.Type == 1 {
		video := form.File["video"][0]
		if video == nil {
			gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_IMAGE_FAIL, nil)
			return
		}
		fileName := upload.GetFileName(video.Filename)
		videoSrc := "/" + setting.MinioSetting.VideoBucketName + "/" + fileName + ".ts"
		preiviewSrc := "/" + setting.MinioSetting.PreviewBucketName + "/" + fileName + ".jpg"

		if !upload.CheckVideoSize(video) {
			gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
			return
		}
		// TODO: 保存到minio,通过消息队列异步压缩视频和生成预览图
		tempSavePath := "runtime/temp/" + fileName
		if err = c.SaveUploadedFile(video, tempSavePath); err != nil {
			log.Logger.Error("保存文件失败", err)
			gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
			return
		}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer os.Remove(tempSavePath)
			defer wg.Done()
			if err = minio.PutVideoAndPreview(tempSavePath, fileName); err != nil {
				log.Logger.Error("保存视频失败", err)
				gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
				// os.Remove(tempSavePath)
				return
			}
			params.VideoUrl = videoSrc
			params.PreviewUrl = preiviewSrc
			log.Logger.Info("保存上传视频", videoSrc)
		}()
		wg.Wait()
		// if err = minio.PutVideoAndPreview(tempSavePath, fileName); err != nil {
		// 	log.Logger.Error("保存视频失败", err)
		// 	gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		// 	os.Remove(tempSavePath)
		// 	return
		// }
		// os.Remove(tempSavePath)
		// params.VideoUrl = videoSrc
		// params.PreviewUrl = preiviewSrc
		// log.Logger.Info("保存上传视频", videoSrc)
	} else {
		// 获取所有图片
		files := form.File["images"]
		if len(files) == 0 {
			gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_IMAGE_FAIL, nil)
			return
		}
		for _, file := range files {
			imageName := upload.GetImageName(file.Filename)
			// savePath := upload.GetImagePath() + imageName
			if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
				gin_http.Response(c, http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
				return
			}
			if err = minio.PutImage(setting.MinioSetting.ImageBucketName, imageName, file); err != nil {
				log.Logger.Error("保存图片失败", err)
				gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
				return
			}
			// if err = c.SaveUploadedFile(file, "."+savePath); err != nil {
			// 	log.Logger.Error("保存文件失败", err)
			// 	gin_http.Response(c, http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
			// 	return
			// }

			params.ImgUrl = append(params.ImgUrl, imageName)
		}
	}

	err = articleService.Add(&params)
	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

	// for _, img_url := range params.ImgUrl {
	//TODO: 异步
	// if params.Type != 1 {
	// 	go func(urls []string) {
	// 		for _, url := range urls {
	// 			_, err = upload.Thumbnailify(url)
	// 			if err != nil {
	// 				log.Logger.Error("图片压缩失败", err)
	// 			}
	// 		}

	// 	}(params.ImgUrl)
	// }

	// }

}

// @Summary 添加文章标签
// @Produce  json
// @Param id path int true "文章ID"
// @Param tag_name formData []string false "标签"
// @Param token header string true "token"
// @Router /api/v1/article/{id}/addTags [post]
func AddArticleTags(c *gin.Context) {

	articleService := service.GetArticleService()
	// httpCode, errCode := articleService.Bind(c)
	// if errCode != e.SUCCESS {
	// 	gin_http.Response(c, httpCode, errCode, nil)
	// 	return
	// }
	articleId, err1 := strconv.Atoi(c.Param("id"))
	params := param.ArticleAddTagsParams{}
	params.ArticleId = articleId

	if err := c.ShouldBind(&params); err != nil || err1 != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	userID, err := articleService.GetUserID(articleId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if !articleService.CheckTokenUid(c, userID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := articleService.AddTags(&params); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_TAG_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}

// @Summary 删除文章标签
// @Produce  json
// @Param id path int true "文章ID"
// @Param tag_name formData []string true "标签"
// @Param token header string true "token"
// @Router /api/v1/article/{id}/deleteTags [delete]
func DeleteArticleTags(c *gin.Context) {

	articleService := service.GetArticleService()
	params := param.ArticleAddTagsParams{}
	articleId, err1 := strconv.Atoi(c.Param("id"))
	params.ArticleId = articleId

	if err := c.ShouldBind(&params); err != nil || err1 != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	//swagger存在bug delete绑定不了数据
	log.Logger.Debug("绑定数据", params)
	userID, err := articleService.GetUserID(articleId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if !articleService.CheckTokenUid(c, userID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if err = articleService.DeleteTags(&params); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_TAG_FAIL, nil)
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
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	log.Logger.Debug("绑定数据", articleId)
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
// @Param title formData string false "title"
// @Param content formData string false "content"
// @Router /api/v1/article/{id} [put]
func UpdateArticle(c *gin.Context) {

	articleService := service.GetArticleService()
	params := param.ArticleUpdateParams{}
	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil || articleId <= 0 {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	params.ArticleId = articleId

	userID, err := articleService.GetUserID(articleId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if !articleService.CheckTokenUid(c, userID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := articleService.Update(&params); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 恢复文章
// @Produce  json
// @Param id path int true "文章ID"
// @Param token header string true "token"
// @Router /api/v1/article/{id}/recover [put]
func RecoverArticle(c *gin.Context) {
	articleService := service.GetArticleService()

	articleId, err := strconv.Atoi(c.Param("id"))
	if err != nil || articleId <= 0 {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	userID, err := articleService.GetUserID(articleId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}
	if !articleService.CheckTokenUid(c, userID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := articleService.Recovery(articleId); err != nil {
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
	param := param.ArticleLikeParams{}
	if err := c.ShouldBind(&param); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	var err error
	param.ArticleId, err = strconv.Atoi(c.Param("id"))
	if err != nil || param.ArticleId <= 0 {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	log.Logger.Debug("绑定数据", param)
	if !articleService.CheckTokenUid(c, uint(param.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}
	if err := articleService.UpdateArticleLike(&param); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_LIKE_ARTICLE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)

}
