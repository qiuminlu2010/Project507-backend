// import { CommentApi } from "@/components/Comment/interface";

/* themeConfigProp */
export interface ThemeConfigProp {
	primary: string;
	isDark: boolean;
}

export interface UserInfo {
	uid: string;
	username: string;
	name?: string;
	avator: string;
}
/* GlobalState */
export interface GlobalState {
	token: string;
	uuid: string;
	uid: number;
	username: string;
	name?: string;
	avatar?: string;
	expireTime: number;
	assemblySize: string;
	language: string;
	themeConfig: ThemeConfigProp;
}

export interface CommentState {
	currentCommentList: any[];
	editor: any;
	scrollBar: any;
	emojiList: Array<any>;
	loadingComments: boolean;
	noMoreComments: boolean;
	gallery: any;
	selectImageUrls: Array<any>;
	clickComment: boolean;
	selectVideoUrl: string;
	selectPreviewUrl: string;
	selectArticleId: number;
	selectItem: any;
	limit: number;
}
/* MenuState */
export interface MenuState {
	isCollapse: boolean;
	menuList: Menu.MenuOptions[];
}

/* TabsState */
export interface TabsState {
	tabsMenuValue: string;
	tabsMenuList: Menu.MenuOptions[];
}

/* AuthState */
export interface AuthState {
	authButtons: {
		[propName: string]: any;
	};
	authRouter: string[];
}
