package service

import (
	"fmt"
	base "qiu/blog/service/base"
	"time"

	// cache "qiu/blog/service/cache"
	"qiu/blog/model"
	article "qiu/blog/service/article"
	msg "qiu/blog/service/msg"
	param "qiu/blog/service/param"
)

type CommentService struct {
	base.BaseService
}

var commentService CommentService

func GetCommentSerivice() *CommentService {
	return &commentService
}

// 1.更新数据库 2.若key article:id:comments 存在...TODO
func (s *CommentService) Add(params *param.CommentAddParams) error {
	if params.ReplyId > 0 {
		return model.AddReply(params.UserId, params.ArticleId, params.ReplyId, params.Content)
	} else {
		err := model.AddComment(params.UserId, params.ArticleId, params.Content)
		if err != nil {
			return err
		}
		// 推送评论消息
		go pushCommentMessage(params.UserId, params.ArticleId, params.Content)
	}
	return nil
}

func (s *CommentService) Like(params *param.LikeCommentParams) error {
	if params.Type == 1 {
		// 推送点赞评论消息
		go pushLikeCommentMessage(params.UserId, params.CommentId)
		return model.AddCommentLike(params.UserId, params.CommentId)
	} else if params.Type == 0 {
		return model.DeleteCommentLike(params.UserId, params.CommentId)
	}
	return nil
}

func (s *CommentService) Get(params *param.CommentGetParams) ([]*model.Comment, error) {
	return model.GetComments(params.ArticleId, params.UserId, params.PageNum, params.PageSize)
}

func (s *CommentService) Delete(commentId int) error {
	return model.DeleteComment(commentId)
}

func (s *CommentService) GetArticleOwnerId(commentId int) (uint, error) {
	return model.GetArticleOwnerIdByCommentId(commentId)
}

func (s *CommentService) GetArticleOwnerIdByReply(commentId int) (uint, error) {
	return model.GetArticleOwnerIdByReplyId(commentId)
}

func pushCommentMessage(userId int, articleId int, commentContent string) {
	articleInfo, err := article.GetArticleCache(articleId)
	if err != nil {
		return
	}
	content := fmt.Sprintf("<p>用户ID:%d 评论了你的动态</p><p>文章ID:%v</p><p><em>%d</em></p>",
		userId, commentContent, articleInfo.ID)
	message := &msg.Message{
		FromUid:   1,
		ToUid:     int(articleInfo.OwnerID),
		Username:  "系统消息",
		Avatar:    "",
		Content:   content,
		CreatedOn: time.Now().Unix(),
	}
	msg.SystemMsg <- message
}

func pushLikeCommentMessage(userId int, commentId int) {
	comment, err := model.GetCommentById(commentId)
	if err != nil {
		return
	}
	content := fmt.Sprintf("<p>用户ID:%d 点赞了你的评论</p><p>评论ID:%d</p><p><em>%v</em></p>",
		userId, commentId, comment.Content)
	message := &msg.Message{
		FromUid:   1,
		ToUid:     int(comment.UserID),
		Username:  "系统消息",
		Avatar:    "",
		Content:   content,
		CreatedOn: time.Now().Unix(),
	}
	msg.SystemMsg <- message
}
