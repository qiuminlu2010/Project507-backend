<template>
	<div class="w-80% h-full mt-20px">
		<!--  -->
		<v3-waterfall
			class="waterfall"
			:list="artileList"
			srcKey="preview_url"
			:gap="16"
			:colWidth="320"
			:distanceToScroll="200"
			:isLoading="loading"
			:isOver="over"
			:isMounted="isMounted"
			scrollBodySelector=".main-box"
			@scrollReachBottom="getNext"
		>
			<template #default="slotProp">
				<el-card class="pic-card" :body-style="{ padding: '0px' }">
					<div
						v-if="slotProp.item.video_url.length !== 0"
						class="right-2 top-2 absolute w-24px h-24px z-20 i-ep-video-camera"
					></div>
					<el-image :src="slotProp.item.preview_url" fit="contain" @click="handlePreview(slotProp.item)" />
					<div class="pic-info">
						<span>{{ slotProp.item.title }}</span>
						<div class="bottom card-header">
							<div class="time">{{ slotProp.item.ctime }}</div>
							<div class="like_content">
								<button class="btn_like" type="button" @click.stop="handleStar(slotProp.item)">
									<div class="svg_wrap">
										<LikeFilled v-if="slotProp.item.star"></LikeFilled>
										<Like v-else></Like>
									</div>
								</button>
								<div class="like_count">
									<span>{{ slotProp.item.like }}</span>
								</div>
							</div>
						</div>
					</div>
				</el-card>
			</template>
		</v3-waterfall>
		<el-dialog v-model="previewVisible" :close-on-click-modal="true" :show-close="false">
			<!-- <div class="dialog-box"> -->
			<div class="absolute right-1 top-1 i-ep-close-bold h-32px w-32px z-20 cursor-pointer" @click="handleCloseDialog">关闭</div>
			<el-row type="flex" align="middle" justify="center" class="w-full h-full box-border">
				<div class="flex w-full h-96vh box-border">
					<VideoShow v-if="store.selectVideoUrl.length != 0" class="left-box flex relative"></VideoShow>
					<ImageShow v-else class="left-box flex relative"></ImageShow>
					<Right :item="articleItem!"></Right>
				</div>
			</el-row>
		</el-dialog>
	</div>
</template>

<script lang="ts" setup>
import { nextTick, onMounted, ref } from "vue";
import "v3-waterfall/style.css";
import { getArtileList, handleStar } from "./api";
import { ViewCard } from "./interface";
import { Like, LikeFilled } from "./icon";
import { CommentStore } from "@/store/modules/comment";
// import { GlobalStore } from "@/store";
// import { getArticleCommentApi } from "@/api/modules/comment";
// import { formatTime } from "../msg/utils";
// import { ElMessage } from "element-plus";
import ImageShow from "./components/ImageShow.vue";
import VideoShow from "./components/VideoShow.vue";
import Right from "./components/Right.vue";

const artileList = ref<ViewCard[]>([]);
const loading = ref(false);
const over = ref(false);
const offset = ref(0);
const limit = 5;
const store = CommentStore();
// const globalStore = GlobalStore();
const previewVisible = ref(false);
const articleItem = ref<ViewCard>();
const isMounted = ref(false);

const fetchList = async (): Promise<void> => {
	loading.value = true;
	let newList = await getArtileList(offset.value, limit);
	if (newList.length === 0) {
		over.value = true;
		return;
	}
	offset.value += limit;
	loading.value = false;
	artileList.value = artileList.value.concat(newList);
	// if (artileList.value.length > 120) over.value = true;
};

onMounted(fetchList);
nextTick(() => {
	isMounted.value = true;
});
let isLoad = false;
const getNext: () => Promise<void> = async (): Promise<void> => {
	if (isLoad) return;
	isLoad = true;
	await fetchList();
	isLoad = false;
};

const handlePreview = async (item: ViewCard) => {
	// previewTitle.value = item.name;
	// previewURL.value = url;
	store.selectItem = item;
	store.selectVideoUrl = item.video_url;
	store.selectPreviewUrl = item.preview_url;
	store.selectArticleId = item.id;
	if (item.image_url.length !== 0) {
		store.selectImageUrls = [
			{
				src: item.image_url
			},
			{
				src: item.image_url
			}
		];
		if (store.gallery !== null) {
			console.log("refresh");
			store.gallery.updateSlides(store.selectImageUrls, store.gallery.index);
		}
	} else {
		store.gallery = null;
	}

	previewVisible.value = true;
	articleItem.value = item;
	store.currentCommentList = [];
	store.handLoadMoreComment();
	// store.noMoreComments = false;
	// store.currentCommentList = [];
	// const res = await getArticleCommentApi({ article_id: item.id, user_id: globalStore.uid });
	// if (res.code == 200) {
	// 	store.currentCommentList = res.data!.datalist;
	// let temp: CommentCard[] = [];
	// res.data?.datalist.forEach(item => {
	// 	temp.push({
	// 		id: item.ID,
	// 		article_id: item.article_id,
	// 		createTime: formatTime(item.created_on),
	// 		content: item.content,
	// 		like: item.like_count,
	// 		is_like: item.is_like,
	// 		username: item.name,
	// 		avatar: item.avatar,
	// 		reply: item.replies
	// 		// user: {
	// 		// 	user_id: item.user_id,
	// 		// 	username: item.username,
	// 		// 	name: item.name,
	// 		// 	avatar: item.avatar
	// 		// }
	// 	});
	// });
	// console.log("获取评论列表", store.currentCommentList);
	// } else {
	// 	ElMessage.error("获取评论列表失败！");
	// }
};
function handleCloseDialog() {
	previewVisible.value = false;
}
</script>

<style lang="scss" scoped>
@import "./index.scss";
.content {
	padding-top: 30px;
	margin-left: 210px;
	background-color: #dee3e7;
}
</style>
