<template>
	<div>
		<ul
			class="pr-10px overflow-auto h-full"
			v-infinite-scroll="commentStore.handLoadMoreComment"
			:infinite-scroll-disabled="disabled"
		>
			<ContentBox v-for="(comment, index) in data" :key="index" :parent-id="comment.id" :data="comment">
				<ReplyBox :parent-id="comment.id" :data="comment.reply" />
			</ContentBox>
			<div v-if="commentStore.loadingComments" class="text-center c-gray-400">加载更多</div>
			<div v-if="commentStore.noMoreComments" class="text-center c-gray-400">评论到底了</div>
		</ul>
	</div>
</template>

<script setup lang="ts">
import { CommentStore } from "@/store/modules/comment";
import { computed } from "vue";
// import { CommentApi } from "./interface";
import ContentBox from "./ContentBox.vue";
import ReplyBox from "./ReplyBox.vue";
const commentStore = CommentStore();

const disabled = computed(() => commentStore.loadingComments || commentStore.noMoreComments);

const data = computed(() => {
	return commentStore.currentCommentList;
});
</script>

<style lang="scss" scoped></style>
