package service

import "qiu/blog/model"

type CommentService struct {
	BaseService
}

var commentService CommentService

func GetCommentSerivice() *CommentService {
	return &commentService
}

// 1.更新数据库 2.若key article:id:comments 存在...TODO
func (s *CommentService) Add(params *CommentAddParams) error {
	return model.AddComment(params.UserId, params.ArticleId, params.Content)
}
func (s *CommentService) Reply(params *CommentAddParams) error {
	return model.AddReply(params.UserId, params.ArticleId, params.ReplyId, params.Content)
}

func (s *CommentService) Get(params *CommentGetParams) ([]*model.Comment, error) {
	return model.GetComments(params.ArticleId, params.PageNum, params.PageSize)
}

// func (s *CommentService) Reply(params *CommentAddParams) error {
// 	return model.AddComment(params.UserId, params.CommentId, params.Content)
// }

func (s *CommentService) Delete(commentId int) error {
	return model.DeleteComment(commentId)
}

func (s *CommentService) DeleteReply(replyId int) error {
	return model.DeleteReply(replyId)
}

func (s *CommentService) GetArticleOwnerId(commentId int) (uint, error) {
	return model.GetArticleOwnerIdByCommentId(commentId)
}

func (s *CommentService) GetArticleOwnerIdByReply(commentId int) (uint, error) {
	return model.GetArticleOwnerIdByReplyId(commentId)
}
