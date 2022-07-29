import { Message } from "@/api/interface/index";
export interface Message {
	from_id: number;
	to_id: number;
	content: string;
	image_url: string;
	created_on: number;
	status: number;
}
// export interface Messages {
// 	sessionId: number;
// 	avatar: string;
// 	name: string;
// 	msg: Message.ResMessages[];
// }

export interface Main {
	sessionList: Array<any>; //User.UserInfo[]
	allSessionList: Array<any>;
	sessionListOffset: number;
	sessionSelectId: number;
	sessionName: string;
	sessionSelected: any;
	sessionAvatar: string;
	// systemSession: Message.ResMessages;
	socket: any;
	onlineStatus: boolean;
	sendInfo: any;
	emojiList: Array<any>;
	chatScrollbar: any;
	chatEditor: any;
	editor: any;
	editorData: any;
	beforePullDown: boolean;
	isPullingDown: boolean;
}
