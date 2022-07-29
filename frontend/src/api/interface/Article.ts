import { ReqPage } from "./index";

// *文章管理
export namespace Article {
	export interface Comments {
		article_id: number;
	}
	export interface Image {
		ID: number;
		url: string;
		thumb_url: string;
	}
	export interface Tag {
		ID: number;
		name: string;
	}
	export interface ReqGetArticleParams extends ReqPage {
		// title: string;
		// content: string;
		user_id?: number;
	}
	export interface ReqAddArticleParams {
		title: string;
		content: string;
		user_id: number;
		tags: string[];
	}
	export interface ReqAddArticleWithImgParams extends ReqAddArticleParams {
		images: File;
	}
	export interface ResArticleList {
		id: number;
		created_on: number;
		owner_id: number;
		owner_name: string;
		owner_username: string;
		owner_avatar: string;
		title: string;
		content: string;
		like_count: number;
		is_like: boolean;
		images: Image[];
		tags: Tag[];
	}
	export interface ResArticleComment {
		ID: number;
		created_on: number;
		user_id: number;
		article_id: number;
		reply_id?: null;
		username: string;
		avatar: any;
		content: string;
		like: number;
		is_like: boolean;
		replies?: ResArticleComment[];
	}
}
