import { defineStore } from "pinia";
import { CommentState } from "../interface";

export const CommentStore = defineStore({
	id: "CommentState",
	state: (): CommentState => ({
		currentCommentList: [],
		editor: null,
		scrollBar: null,
		emojiList: [],
		loadingComments: false,
		noMoreComments: false,
		selectImageUrls: [],
		gallery: null,
		clickComment: false
	}),
	getters: {},
	actions: {
		setCurrentCommentList(CommentList: any) {
			this.currentCommentList = CommentList;
		},
		handleEditorCreated(editor: any) {
			this.editor = editor;
		},
		handleEditorChange() {
			// console.log("加载更多消息");
			console.log("change:", this!.editor.getHtml().trim());
		},
		handLoadMoreComment() {
			console.log("加载更多评论");
			if (this.currentCommentList.length > 10) {
				this.noMoreComments = true;
				return;
			}
			this.loadingComments = true;
			setTimeout(() => {
				this.currentCommentList.push(this.currentCommentList[0]);
				this.loadingComments = false;
			}, 1000);
		}
	}
});
