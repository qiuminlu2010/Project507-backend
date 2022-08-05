<template>
	<div class="right-box flex-1 min-w-30vw">
		<div class="article-box flex flex-col" id="article-box">
			<div class="user">
				<div class="user-avatar">
					<avatar :size="50" :src="'/minio' + item.owner_avatar"></avatar>
					<!-- <el-avatar :size="50" :src="'/base' + item.owner_avatar" /> -->
				</div>
				<div class="userinfo">
					<div class="username">{{ item.owner_username }}</div>
					<div class="time">{{ item.data }}</div>
				</div>
				<div class="follow">
					<el-button class="follow-btn" round>+关注</el-button>
				</div>
			</div>
			<div class="flex flex-col">
				<Fold :unfold="true" line="3">
					<div v-dompurify-html="item.content"></div>
				</Fold>
				<!-- <div class="flex items-center pb-2px pl-0.5rem">{{ item.content }}</div> -->
				<div class="tag-list">
					<el-button class="tag-box" type="primary" v-for="tag in tags" :key="tag" link> #{{ tag.name }} </el-button>
				</div>
			</div>

			<div class="operation-list">
				<div class="op-box">
					<button class="op-btn">
						<el-icon :size="22"><Share /></el-icon>
						<div class="op-name">分享</div>
					</button>
				</div>
				<div class="op-box">
					<button class="op-btn" @click="handleClickComment">
						<el-icon :size="22"><ChatDotSquare /></el-icon>
						<div class="op-name">评论</div>
					</button>
				</div>
				<div class="op-box">
					<button class="op-btn op-like-btn" @click.stop="handleStar(props.item)">
						<el-icon class="op-icon" :size="22"> <LikeFilled v-if="item.star"></LikeFilled> <Like v-else></Like> </el-icon>
						<div class="op-name like_count">{{ item.like }}</div>
					</button>
				</div>
			</div>
		</div>
		<div class="comment-box pl-5% pr-10px flex flex-col flex-1 overflow-auto" id="comment-box">
			<!-- <div class="suggest" v-if="commentStore.currentCommentList.length <= 0">
				<el-icon :size="40" class="hint-icon"><WarningFilled /></el-icon>
				<div class="hint">快来发表你的评论吧</div>
			</div> -->
			<Comment></Comment>
		</div>
	</div>
</template>
<script lang="ts" setup>
import { computed, nextTick, onMounted, onUpdated } from "vue";
import { Share, ChatDotSquare } from "@element-plus/icons-vue";
import { Like, LikeFilled } from "../icon";
import { ViewCard } from "../interface";
// import { GlobalStore } from "@/store";
import Comment from "./Comment.vue";
import Fold from "@/components/Fold";
// import CommentFoot from "./CommentFoot.vue";
// import { formatTime } from "../utils";
// import Comment from "vue-juejin-comment";
import avatar from "vue-avatar/src/avatar.vue";
import { CommentStore } from "@/store/modules/comment";
const commentStore = CommentStore();
// const globalStore = GlobalStore();
// const commentRef = ref();
// const commentList = computed(() => {
// 	return commentStore.currentCommentList;
// });
// commentStore.currentCommentList = [
// 	{
// 		ID: 1,
// 		reply_id: null,
// 		avatar: "https://static.juzicon.com/avatars/avatar-200602130320-HMR2.jpeg?x-oss-process=image/resize,w_100",
// 		username: "落🤍尘",
// 		// // level: 6,
// 		// // address: "来自上海",
// 		content:
// 			"从正月初一开始就每天坚持在写作文，都说二十一天可以养成一个习惯，今天已经是二十四了。真的是习惯了每天打开写作这件事了。今天，想着不在状态就不要写了吧，",
// 		like_count: 2,
// 		is_like: false,
// 		created_on: 1659550495,
// 		replies: []
// 	},
// 	{
// 		ID: 2,
// 		reply_id: null,
// 		avatar: "https://static.juzicon.com/avatars/avatar-20210310192149-vkuj.jpeg?x-oss-process=image/resize,w_100",
// 		username: "碎梦遗忘录",
// 		// level: 5,
// 		// address: "来自北京",
// 		content: "说谎和沉默可以说是现在人类社会里日渐蔓延的两大罪恶。事实上，我们经常说谎，动不动就沉默不语",
// 		like_count: 4,
// 		created_on: 1659550495,
// 		replies: [
// 			{
// 				ID: 11,
// 				reply_id: 2,
// 				avatar:
// 					"https://static.juzicon.com/avatars/avatar-20220310090547-fxvx.jpeg?x-oss-process=image/resize,m_fill,w_100,h_100",
// 				username: "欲知欲忘",
// 				// level: 4,
// 				// address: "来自成都",
// 				content: "沉默，是保护自己。说谎是让自己不被注意，且不被攻击[狗头]",
// 				like_count: 7,
// 				created_on: 1659550495
// 			},
// 			{
// 				ID: 12,
// 				reply_id: 2,
// 				avatar:
// 					"https://static.juzicon.com/avatars/avatar-20220302110828-1hm0.jpeg?x-oss-process=image/resize,m_fill,w_100,h_100",
// 				username: "陵薮市朝",
// 				// level: 3,
// 				// address: "来自杭州",
// 				content: '回复 <span style="color: blue;"">@欲知欲忘:</span> [吃瓜]果真是了',
// 				like_count: 3,
// 				created_on: 1659550495
// 			},
// 			{
// 				ID: 13,
// 				reply_id: 2,
// 				username: "每天至少八杯水",
// 				avatar:
// 					"https://static.juzicon.com/avatars/avatar-20220308235453-v09s.jpeg?x-oss-process=image/resize,m_fill,w_100,h_100",
// 				like_count: 3,
// 				// level: 2,
// 				// address: "来自深圳",
// 				content: '回复 <span style="color: blue;"">@陵薮市朝:</span> 沉默是金[困狗]',
// 				created_on: 1659550495
// 			}
// 		]
// 	},
// 	{
// 		ID: 3,
// 		reply_id: null,
// 		username: "悟二空",
// 		avatar: "https://static.juzicon.com/user/avatar-bf22291e-ea5c-4280-850d-88bc288fcf5d-220408002256-ZBQQ.jpeg",
// 		// level: 1,
// 		// address: "来自苏州",
// 		content: "知道在学校为什么感觉这么困吗？因为学校，是梦开始的地方。[脱单doge]",
// 		like_count: 11,
// 		created_on: 1659550495,
// 		replies: [
// 			{
// 				ID: 14,
// 				reply_id: 3,
// 				avatar:
// 					"https://static.juzicon.com/user/avatar-8b6206c1-b28f-4636-8952-d8d9edec975d-191001105631-MDTM.jpg?x-oss-process=image/resize,m_fill,w_100,h_100",
// 				username: "别扰我清梦*ぁ",
// 				// level: 5,
// 				// address: "来自重庆",
// 				content: "说的对，所以，综上所述，上课睡觉不怪我呀💤",
// 				like_count: 3,
// 				created_on: 1659550495
// 			},
// 			{
// 				ID: 15,
// 				reply_id: 3,
// 				avatar: "https://static.juzicon.com/avatars/avatar-191031205903-I6EP.jpeg?x-oss-process=image/resize,m_fill,w_100,h_100",
// 				username: "三分打铁",
// 				// level: 3,
// 				// address: "来自武汉",
// 				content: " 仔细一想还真有点感伤[大哭2]",
// 				like_count: 3,
// 				created_on: 1659550495
// 			},
// 			{
// 				ID: 16,
// 				avatar:
// 					"https://static.juzicon.com/user/avatar-3cb86a0c-08e7-4305-9ac6-34e0cf4937cc-180320123405-BCV6.jpg?x-oss-process=image/resize,m_fill,w_100,h_100",
// 				reply_id: 3,
// 				username: "Blizzard",
// 				// level: 4,
// 				content: '回复 <span style="color: blue;"">@别扰我清梦*ぁ:</span> 看完打了一个哈切。。。会传染。。。[委屈]',
// 				// // address: "来自广州",
// 				like_count: 9,
// 				created_on: 1659550495
// 			}
// 		]
// 	}
// ];
onMounted(() => {
	console.log("mounted");
});
onUpdated(() => {
	console.log("update");
	// fixHeight();
});
interface ArticleProps {
	item: ViewCard;
}
const props = withDefaults(defineProps<ArticleProps>(), {
	// title: "",
});
// const currentDate = new Date().toLocaleString();
// const commentList = ref<CommentCard[]>([]);
// const tags = ref<Array<string>>(["Tag 1", "Tag 4"]);
const tags = computed(() => {
	return commentStore.selectItem.tags;
});
function handleClickComment() {
	commentStore.clickComment = !commentStore.clickComment;
}
function handleStar(item: ViewCard) {
	const _item = item;
	if (_item.star) {
		_item.like -= 1;
	} else {
		_item.like += 1;
	}
	_item.star = !_item.star;
}
// function handleCommentLike(item: CommentCard) {
// 	if (item.is_like) {
// 		item.like! -= 1;
// 	} else {
// 		item.like! += 1;
// 	}
// 	item.is_like = !item.is_like;
// }
// const commentInput = ref("");
// const scrollHeight = ref("100px");
nextTick(() => {
	// fixHeight();
});

// const fixHeight = () => {
// 	scrollHeight.value = document.getElementById("comment-box")!.offsetHeight - 45 + "px";
// 	console.log("scrollHeight.value", scrollHeight.value);
// };
</script>
<style scoped lang="scss">
@import "../index.scss";
</style>
