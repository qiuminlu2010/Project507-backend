<template>
	<div :class="{ content: true, 'limit-box': isLimit }">
		<v3-waterfall
			:key="forUpdate"
			:list="list"
			class="waterfall"
			srcKey="cover"
			:gap="20"
			:colWidth="320"
			:distanceToScroll="200"
			:isLoading="loading"
			:isOver="over"
			:isMounted="isMounted"
			:scrollBodySelector="isLimit ? '.limit-box' : ''"
			@scrollReachBottom="getNext"
		>
			<template #default="slotProp">
				<el-card class="pic-card" :body-style="{ padding: '0px' }">
					<el-image :src="slotProp.item.image_url" fit="contain" @click="handlePreview(slotProp.item, slotProp.item.image_url)" />
					<div class="pic-info">
						<span>{{ slotProp.item.title }}</span>
						<div class="bottom card-header">
							<div class="time">{{ slotProp.item.data }}</div>
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
					<!-- </div> -->
				</el-card>
			</template>
		</v3-waterfall>
		<!-- <div>
			<button @click="handleLoadMore(20)">加载更多</button>
		</div> -->
		<el-dialog v-model="previewVisible" :close-on-click-modal="true" :show-close="false">
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
import "v3-waterfall/style.css";
import { onMounted, ref } from "vue";
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
const previewVisible = ref(false);
const articleItem = ref<ViewCard>();
// const artileList = ref<ViewCard[]>([]);
const list = ref<unknown[]>([]);
const loading = ref(false);
const over = ref(false);
const isMounted = ref(false);
const isLimit = ref(false);
const forUpdate = ref(0);
// const isLimit = ref(false);
// const fetchList = async (): Promise<void> => {
// 	loading.value = true;
// 	let newList = await getArtileList(5);
// 	// await handleLoadMore(10);
// 	console.log("length", newList.length);
// 	artileList.value = artileList.value.concat(newList);
// 	// list.value = Artile_list.value;
// 	loading.value = false;
// 	if (artileList.value.length > 20) over.value = true;
// };

// nextTick(() => {
// });
onMounted(() => {
	getArtileList(5);
	// handleLoadMore(5);
});
// 加载更多
const fetchList = async (): Promise<void> => {
	loading.value = true;
	let newList = data;
	// let newList = data;
	console.log("length", newList.length);
	list.value = list.value.concat(newList);
	loading.value = false;
	if (list.value.length > 20) over.value = true;
};
// 加载更多
// const fetchList = async (): Promise<void> => {
// 	loading.value = true;
// 	let newList = await getArtileList(10);
// 	console.log("length", newList.length);
// 	artileList.value = artileList.value.concat(newList);
// 	loading.value = false;
// 	if (artileList.value.length > 20) over.value = true;
// };

onMounted(fetchList);
let isLoad = false;
const getNext: () => Promise<void> = async (): Promise<void> => {
	console.log("getNext");
	if (isLoad) return;
	isLoad = true;
	await fetchList;
	isLoad = false;
};
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
		store.gallery.updateSlides(store.selectImageUrls, store.gallery.index);
	}
	previewVisible.value = true;
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
		console.log("获取评论列表", store.currentCommentList);
	} else {
		ElMessage.error("获取评论列表失败！");
	}
};

function handleCloseDialog() {
	previewVisible.value = false;
}
const data = [
	{
		tags: ["vue", "react", "diff"],
		cover: "https://image.gkshi.com/202142113503diff-cover.png",
		outline:
			"在 Vue 2.x 以及 React 中，为了提升应用的性能，都引入了一个概念: Virtual DOM。简单来说，Virtual DOM 就是一棵用 JS 来描述的应用结构树，它与实际 DOM 树相对应。",
		_id: "607fbd150ae17a86a8ca9e4b",
		title: "Vue 2.x 以及 React (<16) diff 算法简单分析",
		created_time: "04-21 2021"
	},
	{
		tags: ["小程序", ""],
		cover: "https://image.gkshi.com/202149153713miniprogram.jpeg",
		outline: "如题，无论是在普通的 Web 开发还是在微信小程序项目的开发中，无疑我们都至少会有测试环境以及正式环境的区别。",
		_id: "6070043d0ae17a86a8ca9e48",
		title: "微信小程序如何自动切换环境接口地址",
		created_time: "04-09 2021"
	},
	{
		tags: ["vue", "webpack", "npm"],
		cover: "https://image.gkshi.com/20211201513vue-cli-element.jpg",
		outline:
			"在工作中，后台项目使用的是 vue 2.x 配合 Element 来实现的。但是由于项目太老，Element 的版本也很早，导致很多功能都有缺失，多数现有需求都要进行重复的扩展。为了方便团队的开发以及提升体验，便想在 Element 的基础上做一些适合现有业务的封装，并且整合发布成 npm 包。",
		_id: "6007d5550ae17a86a8ca9e21",
		title: "vue-cli 如何实现Element-ui二次开发npm包",
		created_time: "01-20 2021"
	},
	{
		tags: ["本站", "typescript", ""],
		cover: "https://image.gkshi.com/2021118224526revision.jpg",
		outline:
			"盼望着，盼望着，春天的，不对，博客改版的脚本近了...从2020年9月份博客关闭以来，历时5个月的改版，现在终于又双叒上线了！😆",
		_id: "60059f0f0ae17a86a8ca9e12",
		title: "这次改版做了什么？",
		created_time: "01-18 2021"
	},
	{
		tags: ["promise", "js"],
		cover: "",
		outline:
			"对于一个前端er来说，Promise 是一个相当熟悉的名词。我们知道它的出现就是为了在一定程度上解决回调地狱的问题。最开始认为 Promise 指的就是 ES6 中的 Promise，实际发现并不是这样，它只是其中的一种实现。更了解发现，其实只要符合 Promise A+ 的规范，都能实现出 Promise，并且能够和其他的实现相通。",
		_id: "5ed61a4af93a5234433bb588",
		title: "写一个自己的符合 Promise A+ 规范的 Promise 类",
		created_time: "06-02 2020"
	},
	{
		tags: ["vue", "es6", "typescript", ""],
		cover: "https://image.gkshi.com/201943095758vue-logo.png",
		outline: "Vue 3.0 都发布 Beta 公测了，我这连 alpha 版的主要源码都没看完...大清亡了啊😂",
		_id: "5ea945dbf93a5234433bb587",
		title: "Vue 3.0 源码学习—effect",
		created_time: "04-29 2020"
	},
	{
		tags: ["jenkins", "docker", "CI/CD"],
		cover: "https://image.gkshi.com/2020421204444jenkins-docker-cover.png",
		outline: "习惯了工作中持续集成部署的流程，现在想要把自己博客的维护也用上jenkins + docker 的模式，来一把酸爽体验！",
		_id: "5e9eecb7f93a5234433bb586",
		title: "前端使用docker和jenkins搭建CI/CD",
		created_time: "04-21 2020"
	},
	{
		tags: ["vue", "es6", "typescript", ""],
		cover: "http://image.gkshi.com/201943095758vue-logo.png",
		outline:
			"离上篇关于 vue 3.0 的源码学习已经过去老久了，这次要学习的是 reactive.ts 文件，很多东西都看了好几遍，觉得算有点理解了才开始写这个。",
		_id: "5e992444f93a5234433bb583",
		title: "Vue 3.0 源码学习—reactive",
		created_time: "04-17 2020"
	},
	{
		tags: ["demo", "test"],
		cover: "http://image.gkshi.com/not-exist-cover.png",
		outline: "如果图片加载失败，会用默认的错误图片展示。",
		_id: "",
		title: "图片裂开",
		created_time: "04-17 2020"
	},
	{
		tags: ["es6", "javascript"],
		cover: "",
		outline: "call、apply 方法在实际开发中还是有用到的，学习了它的相关原理，再自己手写一遍来加深自己的理解。",
		_id: "5e89f3a1f93a5234433bb582",
		title: "手写模拟实现 call、apply 方法",
		created_time: "04-05 2020"
	},
	{
		tags: ["es6", ""],
		cover: "",
		outline: "将几个比较有意思的面试题做一个小小的记录。",
		_id: "5e7c35a796575e7d52d442bd",
		title: "几个面试题总结",
		created_time: "03-26 2020"
	}
];

// function getData(): Promise<unknown[]> {
// 	return new Promise(resolve => {
// 		const ret = [...data, ...data];
// 		setTimeout(() => {
// 			resolve(ret);
// 		}, 2000);
// 	});
// }
// async function getData(): Promise<unknown[]> {
// 	return new Promise(resolve => {
// 		const ret = [...data, ...data];
// 		setTimeout(() => {
// 			resolve(ret);
// 		}, 2000);
// 	});
// }
</script>
<style scoped lang="scss">
@import "./index.scss";
.content {
	padding-top: 30px;
	margin-left: 210px;
	background-color: #dee3e7;
}
.menu {
	position: fixed;
	left: 0;
	width: 200px;
	height: 100vh;
	padding-top: 300px;
	border-right: 1px solid #cccccc;
	p {
		margin: 20px;
		cursor: pointer;
	}
}
.active {
	color: #20c180;
}
.limit-box {
	height: 300px;
	overflow-y: scroll;
}
.waterfall,
.tags {
	width: 80%;
	min-width: 1200px;
	margin: 0 auto;
}
.tags-wrapper {
	padding: 10px 0;
	margin-bottom: 15px;
	background-color: #ffffff;
	border-bottom: 1px solid #c4cdd4;
	.tags {
		.tag {
			display: inline-block;
			width: 100px;
			padding: 5px 0;
			margin: 0 10px 8px 0;
			font-size: 14px;
			color: #576575;
			text-align: center;
			border: 1px solid #e8eaee;
			border-radius: 4px;
			&:hover {
				border-color: #8599ad;
			}
			&:not(.active) {
				cursor: pointer;
			}
			&.active {
				color: #ffffff;
				background-color: #576575;
				border-color: #576575;
			}
		}
	}
}
.brief {
	width: 100%;
	background-color: #ffffff;
	box-sizing: border-box;
	h3,
	> p {
		word-break: break-all;
		word-wrap: break-word;
	}
	h3 {
		padding: 20px 20px 12px;
		font-weight: 400;
		color: #22252a;
		text-align: center;
	}
	> p {
		padding: 0 20px 12px;
		font-size: 13px;
		color: #8a98a8;
	}
}
.outline-bottom {
	display: flex;
	justify-content: space-between;
	padding: 10px;
	background-color: #f9fafb;
	border-top: 1px solid #e7ebef;
	.article-tags,
	time {
		font-size: 12px;
		color: #73828c;
	}
	.article-tags {
		span {
			margin-right: 6px;
			&.tag:hover {
				color: #000000;
				text-decoration: underline;
				cursor: pointer;
			}
		}
	}
}
.list-item {
	box-sizing: border-box;
	.cover-wrapper {
		overflow: hidden;
		background-color: #ffffff;
	}
	&:hover {
		box-shadow: 5px 5px 5px #cccccc;
		.cover {
			animation: scaleImg 0.1s linear forwards;
		}

		@keyframes scaleImg {
			0% {
				transform: scale(1);
			}
			100% {
				transform: scale(1.5);
			}
		}
	}
}
.cover {
	width: 100%;
	vertical-align: middle;
}

/* ipad pro */
@media screen and(max-width: 1024px) {
	header {
		height: 40vh;
		background-size: 100% 100%;
	}
}

@media screen and(max-width: 1000px) {
	.waterfall,
	.tags {
		width: 100%;
		min-width: auto;
		padding: 0 30px;
		box-sizing: border-box;
	}
}

@media screen and(max-width: 576px) {
	header {
		height: 30vh;
		min-height: 300px;
		background-position: 40% 0;
		background-size: 150% 100%;
	}
	.waterfall,
	.tags {
		width: 100%;
		min-width: auto;
		padding: 0;
		box-sizing: border-box;
	}
	.tags {
		width: 90%;
	}
}
</style>
