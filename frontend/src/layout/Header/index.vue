<template>
	<div class="header flex-shrink-0 border-b border-gray-200 shadow-gray-200">
		<div class="mx-24px flex justify-center">
			<el-tooltip content="主页"
				><el-button text size="large" @click="selectHome"
					><div class="i-ep-home-filled cursor-pointer h-24px w-24px"></div></el-button
			></el-tooltip>
			<el-tooltip content="发现"
				><el-button text size="large"><div class="i-ep-compass cursor-pointer h-24px w-24px"></div></el-button
			></el-tooltip>
			<el-tooltip content="消息"
				><el-button text size="large" @click="selectMessage"
					><div class="i-ep-message cursor-pointer h-24px w-24px"></div></el-button
			></el-tooltip>
			<el-tooltip content="我的"
				><el-button text size="large" @click="selectMy"><div class="i-ep-user cursor-pointer h-24px w-24px"></div></el-button
			></el-tooltip>
			<el-tooltip content="发布">
				<el-button text size="large" @click="addNew" bg round
					><div class="i-ep-plus cursor-pointer h-24px w-28px"></div
				></el-button>
			</el-tooltip>
		</div>
		<form class="nav-searchform">
			<el-autocomplete
				v-model="state2"
				:fetch-suggestions="querySearch"
				:trigger-on-focus="false"
				clearable
				class="nav-search-content"
				placeholder="Please Input"
				@select="handleSelectSuggest"
			/>
			<div class="nav-search-btn">
				<el-icon><Search /></el-icon>
			</div>
		</form>
		<div class="ml-20px flex justify-center items-center">
			<!-- Theme -->
			<Avatar v-if="store.token.length !== 0"></Avatar>
			<div v-else>
				<el-button type="primary" round @click="login">登录/注册</el-button>
			</div>
		</div>
		<SwitchDark class="dark ml-20px right-24px absolute"></SwitchDark>
	</div>
	<UploadArticle ref="dialogRef"></UploadArticle>
</template>

<script setup lang="ts">
import { Search } from "@element-plus/icons-vue";
import { ref, onMounted } from "vue";
// import Fullscreen from "./components/Fullscreen.vue";
import Avatar from "./components/Avatar.vue";
import SwitchDark from "@/components/SwitchDark/index.vue";
// import Theme from "./components/Theme.vue";
import UploadArticle from "@/components/UploadArticle/index.vue";
import { GlobalStore } from "@/store";
import { useRouter } from "vue-router";

const store = GlobalStore();
const router = useRouter();
const selectHome = () => {
	router.push({ name: "home" });
};
const selectMessage = () => {
	router.push({ name: "msg" });
};
const selectMy = () => {
	router.push({ name: "mm" });
};
function login() {
	router.push({ name: "login" });
}
interface DialogExpose {
	acceptParams: () => void;
}
const dialogRef = ref<DialogExpose>();
const addNew = () => {
	dialogRef.value!.acceptParams();
};

interface RestaurantItem {
	value: string;
	link: string;
}

const state2 = ref("");

const restaurants = ref<RestaurantItem[]>([]);
const querySearch = (queryString: string, cb: any) => {
	const results = queryString ? restaurants.value.filter(createFilter(queryString)) : restaurants.value;
	// call callback function to return suggestions
	cb(results);
};
const createFilter = (queryString: string) => {
	return (restaurant: RestaurantItem) => {
		return restaurant.value.toLowerCase().indexOf(queryString.toLowerCase()) === 0;
	};
};
const loadAll = () => {
	return [
		{ value: "vue", link: "https://github.com/vuejs/vue" },
		{ value: "element", link: "https://github.com/ElemeFE/element" },
		{ value: "cooking", link: "https://github.com/ElemeFE/cooking" },
		{ value: "mint-ui", link: "https://github.com/ElemeFE/mint-ui" },
		{ value: "vuex", link: "https://github.com/vuejs/vuex" },
		{ value: "vue-router", link: "https://github.com/vuejs/vue-router" },
		{ value: "babel", link: "https://github.com/babel/babel" }
	];
};

const handleSelectSuggest = (item: RestaurantItem) => {
	console.log(item);
};

onMounted(() => {
	restaurants.value = loadAll();
});
</script>

<style scoped lang="scss">
@import "./index.scss";
</style>
