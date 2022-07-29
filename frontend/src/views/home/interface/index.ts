export interface ViewCard {
	src: any;
	id: number;
	title?: string;
	content: string;
	owner_id: number;
	owner_name: string;
	owner_username: string;
	owner_avatar: string;
	star?: boolean;
	like: number;
	image_url: string;
	backgroundColor?: string;
	[attr: string]: any;
}

export interface UserInfo {
	user_id: number;
	avatar: string;
	username: string;
	name: string;
}
export interface CommentCard {
	ID: number;
	created_on: any;
	article_id: number;
	reply_id?: null;
	content: string;
	like?: number;
	is_like?: boolean;
	user: UserInfo;
	replies?: CommentCard[];
}
interface Point {
	rowPerView: number;
}

export type Breakpoints = Record<number, Point>;

export interface WaterfallProps {
	breakpoints: Breakpoints;
	width: number;
	animationDuration: number;
	animationDelay: number;
	animationEffect: string;
	hasAroundGutter: boolean;
	gutter: number;
	list: ViewCard[];
	animationPrefix: string;
}

export interface ItemWidthProps {
	breakpoints: Breakpoints;
	wrapperWidth: number;
	gutter: number;
	hasAroundGutter: boolean;
	initWidth: number;
}

export interface ItemWidthByBreakpointProps extends ItemWidthProps {
	size: number;
}
