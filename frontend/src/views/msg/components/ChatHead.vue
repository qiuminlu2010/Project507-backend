<template>
	<div class="session-header">
		<div class="text-left pl-20px border-box flex h-full w-full items-center">
			<el-avatar shape="circle" class="!block header-avatar" :size="42" fit="cover" :src="'/minio' + store.sessionAvatar" />
			<span class="fw-600">{{ store.sessionName }}</span>
		</div>
	</div>
</template>

<script setup lang="ts">
import { MsgStore } from "..";
const store = MsgStore();
</script>
<style scoped lang="scss">
.session-header {
	// position: absolute;
	z-index: 10;
	display: flex;
	align-items: center;
	width: 100%;
	height: 64px;
	margin-right: 1px;
	background: #ffffff;
	border-bottom: 1px solid #e1e4e8;
	border-top-right-radius: 4px;
	.header-avatar {
		margin-right: 15px;
		border-radius: 50%;
	}
}
</style>
