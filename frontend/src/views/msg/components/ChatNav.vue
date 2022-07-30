<template>
	<div class="lg:w-7% w-30%">
		<el-menu default-active="2" class="h-100%" @select="handleSelect">
			<div align="center" class="mb-20px">
				<el-avatar :src="store.userInfo?.avatar">user</el-avatar>
				<div class="truncate mt-10px text-14px" :title="store.userInfo?.email" style="width: 80%">
					{{ store.userInfo?.email }}
				</div>
			</div>
			<div v-for="item in store.navList" :key="item.index">
				<el-tooltip placement="left">
					<template #content>{{ item.lable }}</template>
					<el-menu-item :index="item.index" class="justify-center fw-bold text-24px">
						<el-badge :value="getUnReadCount(item.index)" :hidden="!getUnReadCount(item.index)" :max="10">
							<div :class="item.icon"></div>
						</el-badge>
					</el-menu-item>
				</el-tooltip>
			</div>
			<el-tooltip placement="left">
				<template #content>退出聊天</template>
				<el-button type="danger" class="mt-40px" size="large" circle @click="logout">
					<i i="ep-switch-button" class="fw-bold"></i>
				</el-button>
			</el-tooltip>
		</el-menu>
	</div>
</template>

<script setup lang="ts"></script>

<style lang="scss" scoped>
.ep-menu {
	.is-active {
		background: var(--ep-color-primary-light-9);
	}
}
</style>
