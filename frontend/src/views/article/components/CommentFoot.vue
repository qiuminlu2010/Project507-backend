<template>
	<div class="comment-input" id="comment-input">
		<!--工具栏-->
		<el-row type="flex" class="mb-10px pl-10px">
			<el-popover placement="top" popper-class="chat-icon-popover" trigger="click">
				<template #reference>
					<div class="text-20px cursor-pointer dark:filter-invert-100">
						<img width="24" height="24" class="block" src="@/assets/images/emoji.svg" />
					</div>
				</template>
				<el-scrollbar class="emoji" height="150px">
					<ul class="m0 p0 pr-2px flex flex-wrap">
						<li
							v-for="(item, row_id) in store.emojiList"
							:key="row_id"
							class="p-5px list-none hover:animate-heart-beat animate-count-animated animate-duration-1s cursor-pointer"
							:title="item"
							@click="selectIcon(item)"
						>
							{{ item }}
						</li>
					</ul>
				</el-scrollbar>
			</el-popover>
			<div class="ml-10px text-20px i-ep-picture-rounded !cursor-pointer">
				<input ref="UploadImageRef" class="opacity-0" name="customerService" type="file" value="" accept="image/*" />
			</div>
			<div class="ml-10px text-20px i-ep-video-camera !cursor-pointer">
				<input ref="referenceUploadVideo" class="opacity-0 hidden" name="customerService" type="file" value="" accept="video/*" />
			</div>
		</el-row>
		<Editor
			style="width: 100%; height: 80px; overflow-y: hidden"
			:defaultConfig="editorConfig"
			mode="simple"
			@onCreated="store.handleEditorCreated"
			@onChange="store.handleEditorChange"
		/>
		<el-button class="submit-btn float-right mt-10px" round>发表</el-button>
	</div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted } from "vue";
import { Editor } from "@wangeditor/editor-for-vue";
import { IEditorConfig } from "@wangeditor/core";
// import { formatTime } from "../utils";
import BScroll from "@better-scroll/core";
import MouseWheel from "@better-scroll/mouse-wheel";
import { CommentStore } from "@/store/modules/comment";
// BScroll.use(PullDown);
BScroll.use(MouseWheel);
const store = CommentStore();
const editorConfig: Partial<IEditorConfig> = {
	placeholder: "请输入内容...",
	maxLength: 200
};

//初始化滚动条
onMounted(() => {
	nextTick(() => {
		// store.chatScrollbar = scrollbarRef.value;
		store.scrollBar = new BScroll(".scroll-wrapper", {
			probeType: 3,
			click: false,
			scrollX: false,
			scrollY: true,
			bounceTime: 800,
			useTransition: false,
			pullDownRefresh: {
				threshold: 70,
				stop: 56
			},
			mouseWheel: {
				speed: 20,
				invert: false,
				easeTime: 300
			}
		});
		store.scrollBar.on("refresh", () => {
			console.log("评论页面切换");
			// store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
		});
		// store.scrollBar.on("pullingDown", loadMoreMessages);
	});
	// console.log("加载聊天页面");
	// store.toBottom();
});
// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
	// const editor = editorRef.value;
	const editor = store.editor;
	if (editor == null) return;
	editor.destroy();
});
// 选择表情
function selectIcon(item: any) {
	console.log("选择表情", item);
	store.editor.value.insertText(item);
}

store.emojiList =
	"😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊 😫 😴 😌 😛 😜 😒 😓 😔 🙁 😖 😑 😶 🙄 😏 😣 😞 😟 😤 😢 😭 💘 👍 👎 👊 👌".split(" ");
</script>
<style scoped lang="scss">
@import "../index.scss";
.session-footer {
	z-index: 10;
	width: 100%;
	border-top: 1px solid #e1e5e8;
}
.comment-input {
	// display: flex;
	padding: 4%;
	border-top: 1px solid #aaaaaa;
	.input-box {
		width: 88%;
		margin-right: 5px;
	}
	.submit-btn {
		width: 64px;
		height: 28px;
		padding: 0 !important;
		font-size: 0.9rem;
		font-weight: 700;
		line-height: 1;
		color: rgba(51 117 185 / 88%);
		cursor: pointer;

		// border: 1px solid;
		border-radius: 2.25rem;

		// box-sizing: border-box;
	}

	// :deep(.w-e-text-container) {
	// color: #999999;
	// background-color: #0d0d0d;
	// }
}
</style>
