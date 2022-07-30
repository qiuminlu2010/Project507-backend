package v1

import (
	"net/http"
	"qiu/backend/pkg/e"
	gin_http "qiu/backend/pkg/http"
	log "qiu/backend/pkg/logging"
	"qiu/backend/pkg/setting"
	service "qiu/backend/service/comment"
	param "qiu/backend/service/param"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 添加评论
// @Produce  json
// @Param user_id formData int true "用户id"
// @Param article_id formData int true "文章id"
// @Param reply_id formData int false "回复评论id"
// @Param content formData string true "内容"
// @Param token header string true "token"
// @Router /api/v1/comment [post]
func AddComment(c *gin.Context) {

	commentService := service.GetCommentSerivice()
	params := param.CommentAddParams{}

	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if !commentService.CheckTokenUid(c, uint(params.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := commentService.Add(&params); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_COMMENT_ADD_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 点赞评论
// @Produce  json
// @Param id path int true "评论id"
// @Param user_id formData int true "用户id"
// @Param type formData int true "操作类型"
// @Param token header string true "token"
// @Router /api/v1/comment/{id}/like [post]
func LikeComment(c *gin.Context) {

	commentService := service.GetCommentSerivice()
	params := param.LikeCommentParams{}

	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil || commentId <= 0 {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if err := c.ShouldBind(&params); err != nil {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	if !commentService.CheckTokenUid(c, uint(params.UserId)) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	params.CommentId = commentId

	if err := commentService.Like(&params); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_COMMENT_LIKE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 获取评论
// @Produce  json
// @Param article_id path int true "文章id"
// @Param user_id query int false "用户id"
// @Param page_num query int false "page_num"
// @Param page_size query int false "page_size"
// @Router /api/v1/comments/{article_id} [get]
func GetComments(c *gin.Context) {
	commentService := service.GetCommentSerivice()
	params := param.CommentGetParams{}

	articleId, err := strconv.Atoi(c.Param("article_id"))
	if err != nil || articleId <= 0 {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	params.ArticleId = articleId

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

	comments, err := commentService.Get(&params)

	if err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_COMMENT_GET_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["datalist"] = comments
	// data["total"] = total
	data["pageNum"] = page
	data["pageSize"] = params.PageSize
	gin_http.Response(c, http.StatusOK, e.SUCCESS, data)
}

// @Summary 回复评论
// @Produce  json
// @Param user_id formData int true "用户id"
// @Param article_id formData int true "文章id"
// @Param id path int true "评论id"
// @Param content formData string true "内容"
// @Param token header string true "token"
// @Router /api/v1/comment/{id}/reply [post]
// func AddReply(c *gin.Context) {

// 	commentService := service.GetCommentSerivice()
// 	params := service.CommentAddParams{}

// 	commentId, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || commentId <= 0 {
// 		log.Logger.Error("绑定错误", err)
// 		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
// 		return
// 	}

// 	if err := c.ShouldBind(&params); err != nil {
// 		log.Logger.Error("绑定错误", err)
// 		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
// 		return
// 	}

// 	if !commentService.CheckTokenUid(c, uint(params.UserId)) {
// 		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
// 		return
// 	}

// 	params.ReplyId = commentId

// 	if err := commentService.Reply(&params); err != nil {
// 		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_COMMENT_REPLY_FAIL, nil)
// 		return
// 	}
// 	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
// }

// @Summary 删除评论
// @Produce  json
// @Param id path int true "评论id"
// @Param token header string true "token"
// @Router /api/v1/comment/{id} [delete]
func DeleteComment(c *gin.Context) {

	commentService := service.GetCommentSerivice()
	commentId, err := strconv.Atoi(c.Param("id"))
	if err != nil || commentId <= 0 {
		log.Logger.Error("绑定错误", err)
		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	articleUserID, err := commentService.GetArticleOwnerId(commentId)
	if err != nil {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
		return
	}

	if !commentService.CheckTokenUid(c, articleUserID) {
		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if err := commentService.Delete(commentId); err != nil {
		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_COMMENT_DELETE_FAIL, nil)
		return
	}
	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
}

// @Summary 删除回复
// @Produce  json
// @Param id path int true "回复id"
// @Param token header string true "token"
// @Router /api/v1/reply/{id} [delete]
// func DeleteReply(c *gin.Context) {

// 	commentService := service.GetCommentSerivice()
// 	replyId, err := strconv.Atoi(c.Param("id"))
// 	if err != nil || replyId <= 0 {
// 		log.Logger.Error("绑定错误", err)
// 		gin_http.Response(c, http.StatusBadRequest, e.INVALID_PARAMS, nil)
// 		return
// 	}
// 	articleUserID, err := commentService.GetArticleOwnerIdByReply(replyId)
// 	if err != nil {
// 		fmt.Println("GetArticleOwnerId", err)
// 		gin_http.Response(c, http.StatusBadRequest, e.ERROR_GET_USERID_FAIL, nil)
// 		return
// 	}

// 	if !commentService.CheckTokenUid(c, articleUserID) {
// 		gin_http.Response(c, http.StatusBadRequest, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
// 		return
// 	}

// 	if err := commentService.DeleteReply(replyId); err != nil {
// 		gin_http.Response(c, http.StatusInternalServerError, e.ERROR_COMMENT_DELETE_FAIL, nil)
// 		return
// 	}
// 	gin_http.Response(c, http.StatusOK, e.SUCCESS, nil)
// }
