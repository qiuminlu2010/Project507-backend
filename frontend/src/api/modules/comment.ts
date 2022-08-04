import { ResPage, Article } from "@/api/interface/index";
import { Comment } from "@/api/interface/comment";
import { API } from "@/api/config/servicePort";
import http from "@/api";
// import { ContentTypeEnum } from "@/enums/httpEnum";

//*获取文章评论
export const getArticleCommentApi = (params: Comment.ReqGetParams) => {
	return http.get<ResPage<Article.ResArticleComment>>(API + `/comments/` + params.article_id, params);
};

//*获取评论点赞
export const getCommentLike = (params: { id: number; user_id: number; type: number; token: string }) => {
	return http.post(API + `/comment/` + params.id + "/like", params);
};

//*添加评论
export const addComment = (params: Comment.ReqAddParams) => {
	return http.post(API + `/comment`, params);
};

// * 添加评论点赞
export const addCommentLike = (params: { id: number; user_id: number; type: number }) => {
	return http.post(API + `/comment/` + params.id + "/like", params);
};
