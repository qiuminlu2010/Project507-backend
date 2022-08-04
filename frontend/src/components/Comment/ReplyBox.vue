<template>
	<div v-if="data && data!.length > 0" class="reply-box">
		<div class="comment-list pl-10px">
			<ContentBox v-for="(reply, index) in data" :key="index" :parent-id="parentId" :data="reply" small />
			<!-- <div class="fetch-more">
        <span class="fetch-more-comment">
          查看更多回复<svg width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg" class=""><path data-v-d6f79dbc="" fill-rule="evenodd" clip-rule="evenodd" d="M5.99976 7.93206L10.0656 3.86619C10.1633 3.76856 10.3215 3.76856 10.4192 3.86619L10.7727 4.21975C10.8704 4.31738 10.8704 4.47567 10.7727 4.5733L6.35331 8.99272C6.15805 9.18798 5.84147 9.18798 5.6462 8.99272L1.22679 4.5733C1.12916 4.47567 1.12916 4.31738 1.22679 4.21975L1.58034 3.86619C1.67797 3.76856 1.83626 3.76856 1.93389 3.86619L5.99976 7.93206Z"></path></svg>
        </span>
      </div> -->
		</div>
	</div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import ContentBox from "./ContentBox.vue";
import type { CommentApi } from "./interface";
// import avatar from "vue-avatar/src/avatar.vue";
interface Props {
	data?: CommentApi[] | null;
	parentId: number;
}

const props = defineProps<Props>();

const data = computed(() => {
	if (!props.data) {
		return {
			total: 0,
			list: []
		};
	}
	return props.data;
});
</script>

<style lang="scss" scoped>
.reply-box {
	margin-top: 16px;
	.comment-list {
		// padding: 16px;
		background: rgb(247 248 250 / 70%);
		border-radius: 4px;
	}
	.fetch-more {
		margin-top: 24px;
		margin-left: 36px;
		.fetch-more-comment {
			font-size: 14px;
			line-height: 22px;
			color: #252933;
			cursor: pointer;
			box-sizing: border-box;
			svg {
				fill: currentcolor;
				margin-left: 4px;
			}
		}
	}
}
</style>
