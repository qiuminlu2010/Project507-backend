import { ResPage, Article } from "@/api/interface/index";
import { API } from "@/api/config/servicePort";
import http from "@/api";
import { ContentTypeEnum } from "@/enums/httpEnum";

/**
 * @name 文章管理模块
 */

// * 获取文章列表
export const getArticleListApi = (params: Article.ReqGetArticleParams) => {
	return http.get<ResPage<Article.ResArticleList>>(API + `/article`, params);
};

// * 切换文章状态
export const changeArticleStatus = (params: { id: number; state: number }) => {
	return http.put(API + `/article/update/` + params.id + `/state`, params);
};

// * 添加文章点赞
export const addArticleLike = (params: { id: number; user_id: number; type: number }) => {
	return http.post(API + `/article/` + params.id + "/like", params);
};

// * 新增文章
export const addArticleApi = (params: any) => {
	return http.post(API + `/article`, params, { headers: { "Content-Type": ContentTypeEnum.FORM_DATA } });
};

// * 更新文章
export const editArticle = (params: { id: number }) => {
	return http.post(API + `/article/update/` + params.id);
};

// * 删除文章
export const deleteArticle = (params: { id: number }) => {
	return http.delete(API + `/article/delete/` + params.id);
};

// * 导出用户数据
export const exportArticleInfo = (params: Article.ReqGetArticleParams) => {
	return http.post<BlobPart>(API + `/user/export`, params, { responseType: "blob" });
};
