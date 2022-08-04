<template>
	<div>
		<ul
			class="overflow-auto h-full w-full"
			v-infinite-scroll="commentStore.handLoadMoreComment"
			:infinite-scroll-disabled="disabled"
		>
			<ContentBox v-for="(comment, index) in data" :key="index" :parent-id="comment.ID" :data="comment">
				<ReplyBox :parent-id="comment.ID" :data="comment.replies" />
			</ContentBox>
			<div v-if="commentStore.loadingComments" class="text-center c-gray-400">加载更多</div>
			<div v-if="commentStore.noMoreComments" class="text-center c-gray-400">评论到底了</div>
		</ul>
	</div>
</template>

<script setup lang="ts">
import { CommentStore } from "@/store/modules/comment";
import { computed } from "vue";
import ContentBox from "./ContentBox.vue";
import ReplyBox from "./ReplyBox.vue";
import { CommentApi } from "./interface";
const commentStore = CommentStore();

const disabled = computed(() => commentStore.loadingComments || commentStore.noMoreComments);

const data = computed<CommentApi[]>(() => {
	return commentStore.currentCommentList;
});
</script>

<style lang="scss" scoped></style>
