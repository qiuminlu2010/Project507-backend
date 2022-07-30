<template>
	<div class="home_box">
		<Waterfall
			:list="Artile_list"
			:row-key="options.rowKey"
			:gutter="options.gutter"
			:has-around-gutter="options.hasAroundGutter"
			:width="options.width"
			:breakpoints="options.breakpoints"
			:img-selector="options.imgSelector"
			:animation-effect="options.animationEffect"
			:animation-duration="options.animationDuration"
			:animation-delay="options.animationDelay"
			:lazyload="options.lazyload"
			:load-props="options.loadProps"
		>
			<template #item="{ item, url }">
				<el-card class="pic-card" :body-style="{ padding: '0px' }">
					<LazyImg :url="url" @click="handlePreview(item, url)" />
					<div class="pic-info">
						<span>{{ item.title }}</span>
						<div class="bottom card-header">
							<div class="time">{{ item.data }}</div>
							<div class="like_content">
								<button class="btn_like" type="button" @click.stop="handleStar(item)">
									<div class="svg_wrap">
										<LikeFilled v-if="item.star"></LikeFilled>
										<Like v-else></Like>
									</div>
								</button>
								<div class="like_count">
									<span>{{ item.like }}</span>
								</div>
							</div>
						</div>
					</div>
					<!-- </div> -->
				</el-card>
			</template>
		</Waterfall>

		<div>
			<button @click="handleLoadMore(20)">加载更多</button>
		</div>
		<el-dialog v-model="previewVisible" :close-on-click-modal="true">
			<!-- <div class="dialog-box"> -->
			<div class="absolute right-1 top-1 i-ep-close-bold h-32px w-32px z-20 cursor-pointer" @click="handleCloseDialog">关闭</div>
			<el-row type="flex" align="middle" justify="center" class="w-full h-full box-border">
				<div class="flex w-full h-96vh box-border">
					<ImageShow class="left-box flex relative"></ImageShow>
					<Right :item="articleItem!"></Right>
				</div>
			</el-row>
		</el-dialog>
	</div>
</template>

<script lang="ts" setup>
import { LazyImg, Waterfall } from "vue-waterfall-plugin-next";
import "vue-waterfall-plugin-next/style.css";
import { onMounted, reactive, ref } from "vue";
import loading from "./assets/loading.png";
import error from "./assets/error.png";
import { getArtileList, handleStar } from "./api";
import { CommentCard, ViewCard } from "./interface";
import { Like, LikeFilled } from "./icon";
import ImageShow from "./components/ImageShow.vue";
import Right from "./components/Right.vue";
import { getArticleCommentApi } from "@/api/modules/comment";
import { ElMessage } from "element-plus";
import { GlobalStore } from "@/store";
import { formatTime } from "./utils";
import { CommentStore } from "@/store/modules/comment";
const store = CommentStore();
const globalStore = GlobalStore();
const { Artile_list, options, handleLoadMore } = useWaterfall();
const previewVisible = ref(false);
// const previewTitle = ref<string | undefined>("");
// const previewURL = ref("");
// const srcList = ref<string[]>([]);
const articleItem = ref<ViewCard>();

const handlePreview = async (item: ViewCard, url: string) => {
	// previewTitle.value = item.name;
	// previewURL.value = url;
	store.selectImageUrls = [
		{
			src: url
		},
		{
			src: url
		}
	];
	if (store.gallery !== null) {
		console.log("refresh");
		// store.gallery.goToNextSlide();
		store.gallery.updateSlides(store.selectImageUrls, store.gallery.index);
		// store.gallery.updateSlides([{ src: url }], store.gallery.index);
	}
	// store.selectImageUrls = [url, url];
	previewVisible.value = true;
	// srcList.value = [url, url];
	articleItem.value = item;
	const res = await getArticleCommentApi({ article_id: item.id, user_id: globalStore.uid });
	if (res.code == 200) {
		let temp: CommentCard[] = [];
		res.data?.datalist.forEach(item => {
			temp.push({
				ID: item.ID,
				article_id: item.article_id,
				created_on: formatTime(item.created_on),
				content: item.content,
				like: item.like_count,
				is_like: item.is_like,
				user: {
					user_id: item.user_id,
					username: item.username,
					name: item.name,
					avatar: item.avatar
				}
			});
		});
		// store.setCurrentCommentList(temp);
		console.log("获取评论列表", store.currentCommentList);
	} else {
		ElMessage.error("获取评论列表失败！");
	}
};

function useWaterfall() {
	const Artile_list = ref<ViewCard[]>([]);
	const options = reactive({
		// 唯一key值
		rowKey: "id",
		// 卡片之间的间隙
		gutter: 20,
		// 是否有周围的gutter
		hasAroundGutter: true,
		// 卡片在PC上的宽度
		width: 340,
		// 自定义行显示个数，主要用于对移动端的适配
		breakpoints: {
			1200: {
				// 当屏幕宽度小于等于1200
				rowPerView: 3
			},
			800: {
				// 当屏幕宽度小于等于800
				rowPerView: 3
			},
			500: {
				// 当屏幕宽度小于等于500
				rowPerView: 2
			}
		},
		// 动画效果
		animationEffect: "animate__fadeInUp",
		// 动画时间
		animationDuration: 1000,
		// 动画延迟
		animationDelay: 300,

		// imgSelector
		imgSelector: "src.original",
		// 加载配置
		loadProps: {
			loading,
			error
		},
		// 是否懒加载
		lazyload: true
	});

	onMounted(() => {
		handleLoadMore(20);
	});

	// 加载更多
	async function handleLoadMore(cnt: number) {
		Artile_list.value.push(...(await getArtileList(cnt)));
	}
	return {
		Artile_list,
		options,
		handleLoadMore
	};
}
function handleCloseDialog() {
	previewVisible.value = false;
}
</script>
<style scoped lang="scss">
@import "./index.scss";
</style>
