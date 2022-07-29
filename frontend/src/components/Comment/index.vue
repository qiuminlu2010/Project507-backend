<template>
	<div class="u-comment box-border">
		<div class="comment-form">
			<div class="header">
				<span class="header-title">评论</span>
			</div>
			<div class="content">
				<div class="mr-16px">
					<avatar :size="40" src="https://cube.elemecdn.com/e/fd/0fc7d20532fdaf769a25683617711png.png"></avatar>
				</div>
				<CommentBox placeholder="输入评论（Enter换行，Ctrl + Enter发送）" content-btn="发表评论" />
			</div>
		</div>
		<!-- <div class="hot-list"></div> -->
		<div class="comment-list-wrapper">
			<slot name="list-title">
				<div class="title">全部评论</div>
			</slot>
			<CommentList />
		</div>
	</div>
</template>

<script setup lang="ts">
import { provide } from "vue";
import CommentBox from "./CommentBox.vue";
import CommentList from "./CommentList.vue";
import avatar from "vue-avatar/src/avatar.vue";
import {
	// InjectionUserApi,
	// CommentSubmitParam,
	// UserApi,
	// EmojiApi,
	// CommentApi,
	// InjectionCommentFun,
	// InjectionEmojiApi,
	InjectionLikeFun,
	InjectionLinkFun
} from "./interface";

// interface Props {
// comments: CommentApi[];
// emoji: EmojiApi;
// user: UserApi;
// }

// const props = defineProps<Props>();

const emit = defineEmits<{
	// (e: "submit", obj: CommentSubmitParam): void;
	(e: "like", id: number): void;
	(e: "link", url: string): void;
}>();

// const submit = (obj: CommentSubmitParam) => {
// 	emit("submit", obj);
// };

const like = (id: number) => {
	emit("like", id);
};

// provide(InjectionCommentFun, submit);
// provide(InjectionEmojiApi, props.emoji);
// provide(InjectionUserApi, props.user);
provide(InjectionLikeFun, like);
provide(InjectionLinkFun, (url: string) => emit("link", url));
</script>

<style lang="scss" scoped>
.u-comment {
	padding: 0 2.5rem;
	padding-bottom: 2rem;

	// margin-top: 1.5rem;
	background-color: #ffffff;
	border-radius: 4px;

	// 设置的边框和内边距的值是包含在width内
	box-sizing: border-box;
	.comment-form {
		padding-top: 2rem;
		.header {
			.header-title {
				font-size: 18px;
				font-weight: 600;
				line-height: 30px;
				color: #000000;
			}
		}
		.content {
			display: flex;
			margin-top: 2rem;
			.el-avatar {
				margin-right: 16px;
			}
		}
	}
	.comment-list-wrapper {
		padding: 40px 0 16px;
		.title {
			position: relative;
			display: flex;
			justify-content: space-between;
			align-items: center;
			width: 100%;
			padding-bottom: 8px;
			font-size: 18px;
			font-weight: 600;
			line-height: 30px;
			color: #252933;
		}
	}
}
</style>
