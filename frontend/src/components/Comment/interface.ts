import { InjectionKey } from "vue";

export interface CommentSubmitParam {
	clear: () => void;
	content: string;
	parentId?: number;
}

export interface CommentApi {
	ID: number;
	article_id?: number;
	reply_id: number | null;
	avatar: string;
	username: string;
	name?: string;
	user_id?: number;
	// level: number;
	// address: string;
	content: string;
	like_count: number;
	created_on: number;
	is_like?: boolean;
	replies?: CommentApi[] | null;
	// reply?: ReplyApi | null;
}

export interface UserApi {
	id: number;
	username: string;
	avatar: string;
	likes: number[];
}

export interface ReplyApi {
	total: number;
	list: CommentApi[];
}

export interface Emoji {
	[key: string]: string;
}
export interface EmojiApi {
	faceList: string[];
	emojiList: Emoji[];
	allEmojiList: Emoji;
}

export const InjectionCommentFun: InjectionKey<(obj: CommentSubmitParam) => void> = Symbol();
export const InjectionEmojiApi: InjectionKey<EmojiApi> = Symbol();
export const InjectionUserApi: InjectionKey<UserApi> = Symbol();
export const InjectionLikeFun: InjectionKey<(id: number) => void> = Symbol();
export const InjectionLinkFun: InjectionKey<() => void> = Symbol();
