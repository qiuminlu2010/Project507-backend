// * 请求响应参数(不包含data)
export interface Result {
	code: number;
	msg: string;
}

// * 请求响应参数(包含data)
export interface ResultData<T = any> extends Result {
	data?: T;
}

// * 分页响应参数
export interface ResPage<T> {
	datalist: T[];
	page_num: number;
	page_size: number;
	total?: number;
}

// * 分页请求参数
export interface ReqPage {
	page_num: number;
	page_size: number;
}
export namespace Page {
	export interface Request {
		offset: number;
		limit: number;
	}
	export interface Response<T> {
		datalist: T[];
		offset: number;
		limit: number;
		total?: number;
	}
}
// * 登录
export namespace Login {
	export interface ReqLoginForm {
		username: string;
		password: string;
	}
	export interface ResLogin {
		uid: number;
		username: string;
		name: string;
		avator: string;
		fan_num: string;
		follow_num: string;
		uuid: string;
		token: string;
		expire_time: number;
	}
	export interface ResAuthButtons {
		[propName: string]: any;
	}
}

// * 用户管理
export namespace User {
	export interface UserInfo {
		id: number;
		username: string;
		name: string;
		avatar: string;
	}
	export interface ReqGetUserParams extends ReqPage {
		username: string;
		gender: number;
		idCard: string;
		email: string;
		address: string;
		createTime: string[];
		status: number;
	}
	export interface ResUserList {
		ID: number;
		username: string;
		password: string;
		// gender: string;
		// age: number;
		student_id: string;
		email: string;
		address: string;
		created_on: string;
		modified_on: string;
		avatar: string;
		state: number;
		children?: ResUserList[];
	}
}

// *标签管理
export namespace Tag {
	export interface ReqTag extends ReqPage {
		name: string;
	}
	export interface ResTag {
		ID: number;
		name: string;
		state: number;
	}
}
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
		uid?: number;
		// tags: string[];
		// images: File[];
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
		state: number;
		// modified_on: string;
	}
	export interface ResArticleComment {
		ID: number;
		created_on: number;
		user_id: number;
		article_id: number;
		reply_id?: null;
		username: string;
		name: string;
		avatar: any;
		content: string;
		like_count: number;
		is_like: boolean;
		replies?: ResArticleComment[];
	}
}

// *消息管理
export namespace Message {
	export interface MessageInfo {
		from_uid: number;
		to_uid: number;
		content: string;
		image_url: string;
		video_url?: string;
		created_on: number;
		status: number;
		readed?: boolean;
	}
	export interface ReqGetParams extends Page.Request {
		from_uid: number;
		to_uid: number;
	}
	export interface ReqGetSessionParams extends ReqPage {
		uid: number;
	}
	export interface ResMessages {
		messages: MessageInfo[];
		userInfo: User.UserInfo;
	}
	export interface SessionInfo extends User.UserInfo {
		unread: number;
		messages: MessageInfo[];
	}
}

export interface ImageInfo {
	image_url: string;
}

export interface VideoInfo {
	video_url: string;
}
