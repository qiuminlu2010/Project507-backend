import { getArticleCommentApi } from "@/api/modules/comment";
import { ElMessage } from "element-plus";
import { defineStore } from "pinia";
import { GlobalStore } from "..";
import { CommentState } from "../interface";
import { Comment } from "@/api/interface/comment";
const globalStore = GlobalStore();
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
		clickComment: false,
		selectVideoUrl: "",
		selectPreviewUrl: "",
		selectArticleId: 0,
		selectItem: null,
		limit: 5
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
			if (!this.selectItem) return;
			if (this.noMoreComments) {
				// this.noMoreComments = true;
				return;
			}
			this.loadingComments = true;
			setTimeout(async () => {
				let params: Comment.ReqGetParams = {
					article_id: this.selectItem.id,
					user_id: globalStore.uid,
					offset: 0,
					limit: this.limit
				};
				const res = await getArticleCommentApi(params);
				if (res.code == 200) {
					let newList = res.data!.datalist;
					this.currentCommentList = this.currentCommentList.concat(newList);
					console.log("获取评论列表", this.currentCommentList);
					if (newList.length < this.limit) this.noMoreComments = true;
				} else {
					ElMessage.error("获取评论列表失败！");
					this.noMoreComments = true;
				}
				// this.currentCommentList.push(this.currentCommentList[0]);
				this.loadingComments = false;
			}, 1000);
		}
	}
});
