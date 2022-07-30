<template>
	<el-container class="justify-center flex flex-col w-full h-full">
		<el-header><Header></Header></el-header>
		<el-main class="main-box justify-center flex w-full box-border flex-1 h-0 overflow-y-auto">
			<router-view v-slot="{ Component, route }">
				<transition appear name="fade-transform" mode="out-in">
					<keep-alive>
						<component :is="Component" :key="route.path"></component>
					</keep-alive>
				</transition>
			</router-view>
		</el-main>
		<el-footer><Footer></Footer></el-footer>
	</el-container>
</template>

<script setup lang="ts">
import Header from "./Header/index.vue";
// import Tabs from "./Tabs/index.vue";
import Footer from "./Footer/index.vue";
// import cacheRouter from "@/routers/cacheRouter";
import { onMounted } from "vue";
// import { getAuthButtons } from "@/api/modules/login";
// import { AuthStore } from "@/store/modules/auth";

// const authStore = AuthStore();

onMounted(async () => {
	// 获取按钮权限列表
	// const res = await getAuthButtons();
	// res.data && authStore.setAuthButtons(res.data);
});
</script>

<style scoped lang="scss">
@import "./index.scss";
</style>
