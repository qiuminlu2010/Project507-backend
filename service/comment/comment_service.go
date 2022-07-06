package service

import (
	base "qiu/blog/service/base"
	// cache "qiu/blog/service/cache"
	"qiu/blog/model"
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
		return model.AddComment(params.UserId, params.ArticleId, params.Content)
	}
}

func (s *CommentService) Like(params *param.LikeCommentParams) error {
	if params.Type == 1 {
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
