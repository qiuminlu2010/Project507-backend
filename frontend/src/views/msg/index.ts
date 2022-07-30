import { defineStore } from "pinia";
import { Main } from "./interface";
export const MsgStore = defineStore({
	id: "msg",
	state: (): Main => {
		return {
			sessionList: [],
			allSessionList: [],
			sessionSelectId: 0,
			sessionName: "",
			sessionAvatar: "",
			sessionListOffset: 0,
			sessionSelected: null,
			socket: null,
			onlineStatus: false,
			sendInfo: null,
			emojiList: [],
			chatScrollbar: null,
			chatEditor: null,
			editor: null,
			editorData: null,
			beforePullDown: true,
			isPullingDown: false
		};
	},
	actions: {
		// 本地新增信息记录
		sendLocal() {},
		// 初始化编辑器

		// 设置会话窗口到达底部
		toBottom() {
			// this.chatScrollbar?.setScrollTop(9999);
			// const timer = setTimeout(() => {
			// 	// this.chatScrollbar?.setScrollTop(9999);
			// 	this.chatScrollbar.scrollTo(0, this.chatScrollbar.maxScrollY);
			// 	clearTimeout(timer);
			// }, 100);

			setTimeout(() => {
				// store.chatScrollbar.refresh();
				this.chatScrollbar.refresh();
				this.chatScrollbar.scrollTo(0, this.chatScrollbar.maxScrollY);
			}, 20);
		},

		insertText(data: string) {
			this.editor.focus();
			setTimeout(() => {
				console.log("插入", this.editor == null);
				this.editor.insertText(data);
			}, 100);
		}
	}
});
