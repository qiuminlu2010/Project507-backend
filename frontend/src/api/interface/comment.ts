import { Page } from ".";

export namespace Comment {
	export interface ReqAddParams {
		user_id: number;
		article_id: number;
		reply_id?: number;
		content: string;
	}
	export interface ReqGetParams extends Page.Request {
		user_id: number;
		article_id: number;
	}
}
