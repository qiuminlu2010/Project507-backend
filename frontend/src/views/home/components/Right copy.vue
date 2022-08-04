<template>
	<div class="right-box">
		<div class="article-box" id="article-box">
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
			<div class="article-content">{{ item.content }}</div>
			<div class="tag-list">
				<el-button class="tag-box" type="primary" v-for="tag in tags" :key="tag" link>
					{{ tag }}
				</el-button>
			</div>
			<div class="operation-list">
				<div class="op-box">
					<button class="op-btn">
						<el-icon :size="22"><Share /></el-icon>
						<div class="op-name">分享</div>
					</button>
				</div>
				<div class="op-box">
					<button class="op-btn">
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
		<div class="comment-box" id="comment-box">
			<div class="header">评论</div>
			<div class="divider"></div>
			<div class="suggest" v-if="commentList.length <= 0">
				<el-icon :size="40" class="hint-icon"><WarningFilled /></el-icon>
				<div class="hint">快来发表你的评论吧</div>
			</div>
			<div class="comment-list max-h-100%" v-else>
				<Comment></Comment>
				<!-- <div ref="commentRef" class="comment-wrap">
					<Comment
						:data="data1"
						:user="currentUser"
						:before-submit="submitComment"
						:before-like="likeComment"
						:before-delete="deleteComment"
						:upload-img="uploadImg"
					></Comment>
				</div> -->
				<!-- <div class="comment-item" v-for="comment in commentList" :key="comment.ID">
					<div class="comment-left">
						<el-avatar :src="comment.avatar" />
					</div>
					<div class="comment-right">
						<div class="content">
							<span class="username">{{ comment.username }} </span> {{ comment.content }}
						</div>
						<div class="bottom">
							<div class="time">{{ formatTime(comment.created_on) }}</div>
							<div class="reply">
								<button class="cm-btn">
									<div>回复</div>
								</button>
							</div>
							<button class="cm-btn cm-like-btn" @click.stop="handleCommentLike(comment)">
								<el-icon :size="13"><ThumbFilled v-if="comment.is_like"></ThumbFilled> <Thumb v-else></Thumb></el-icon>
								<div class="like_count">{{ comment.like }}</div>
							</button>
						</div>
					</div>
				</div> -->
			</div>
		</div>
		<CommentFoot></CommentFoot>
	</div>
</template>
<script lang="ts" setup>
import { computed, nextTick, onMounted, onUpdated, ref } from "vue";
import { Share, ChatDotSquare, WarningFilled } from "@element-plus/icons-vue";
import { Like, LikeFilled } from "../icon";
import { ViewCard } from "../interface";
// import { GlobalStore } from "@/store";
import Comment from "@/components/Comment/index.vue";
import CommentFoot from "./CommentFoot.vue";
// import { formatTime } from "../utils";
// import Comment from "vue-juejin-comment";
import avatar from "vue-avatar/src/avatar.vue";
import { CommentStore } from "@/store/modules/comment";
const commentStore = CommentStore();
// const globalStore = GlobalStore();
// const commentRef = ref();
const commentList = computed(() => {
	return commentStore.currentCommentList;
});
commentStore.currentCommentList = [
	{
		id: 1,
		parentId: null,
		avatar: "https://static.juzicon.com/avatars/avatar-200602130320-HMR2.jpeg?x-oss-process=image/resize,w_100",
		username: "落🤍尘",
		level: 6,
		address: "来自上海",
		content:
			"缘生缘灭，缘起缘落，我在看别人的故事，别人何尝不是在看我的故事?别人在演绎人生，我又何尝不是在这场戏里?谁的眼神沧桑了谁?我的眼神，只是沧桑了自己[喝酒]",
		like: 2,
		createTime: "1分钟前",
		reply: null
	},
	{
		id: 2,
		parentId: null,
		avatar: "https://static.juzicon.com/avatars/avatar-20210310192149-vkuj.jpeg?x-oss-process=image/resize,w_100",
		username: "碎梦遗忘录",
		level: 5,
		address: "来自北京",
		content: "说谎和沉默可以说是现在人类社会里日渐蔓延的两大罪恶。事实上，我们经常说谎，动不动就沉默不语",
		like: 4,
		createTime: "2分钟前",
		reply: {
			total: 3,
			list: [
				{
					id: 11,
					parentId: 2,
					avatar:
						"https://static.juzicon.com/avatars/avatar-20220310090547-fxvx.jpeg?x-oss-process=image/resize,m_fill,w_100,h_100",
					username: "欲知欲忘",
					level: 4,
					address: "来自成都",
					content: "沉默，是保护自己。说谎是让自己不被注意，且不被攻击[狗头]",
					like: 7,
					createTime: "1天前"
				},
				{
					id: 12,
					parentId: 2,
					avatar:
						"https://static.juzicon.com/avatars/avatar-20220302110828-1hm0.jpeg?x-oss-process=image/resize,m_fill,w_100,h_100",
					username: "陵薮市朝",
					level: 3,
					address: "来自杭州",
					content: '回复 <span style="color: blue;"">@欲知欲忘:</span> [吃瓜]果真是了',
					like: 3,
					createTime: "2天前"
				},
				{
					id: 13,
					parentId: 2,
					username: "每天至少八杯水",
					avatar:
						"https://static.juzicon.com/avatars/avatar-20220308235453-v09s.jpeg?x-oss-process=image/resize,m_fill,w_100,h_100",
					like: 3,
					level: 2,
					address: "来自深圳",
					content: '回复 <span style="color: blue;"">@陵薮市朝:</span> 沉默是金[困狗]',
					createTime: "5天前"
				}
			]
		}
	},
	{
		id: 3,
		parentId: null,
		username: "悟二空",
		avatar: "https://static.juzicon.com/user/avatar-bf22291e-ea5c-4280-850d-88bc288fcf5d-220408002256-ZBQQ.jpeg",
		level: 1,
		address: "来自苏州",
		content: "知道在学校为什么感觉这么困吗？因为学校，是梦开始的地方。[脱单doge]",
		like: 11,
		createTime: "1天前",
		reply: {
			total: 3,
			list: [
				{
					id: 14,
					parentId: 3,
					avatar:
						"https://static.juzicon.com/user/avatar-8b6206c1-b28f-4636-8952-d8d9edec975d-191001105631-MDTM.jpg?x-oss-process=image/resize,m_fill,w_100,h_100",
					username: "别扰我清梦*ぁ",
					level: 5,
					address: "来自重庆",
					content: "说的对，所以，综上所述，上课睡觉不怪我呀💤",
					like: 3,
					createTime: "1分钟前"
				},
				{
					id: 15,
					parentId: 3,
					avatar:
						"https://static.juzicon.com/avatars/avatar-191031205903-I6EP.jpeg?x-oss-process=image/resize,m_fill,w_100,h_100",
					username: "三分打铁",
					level: 3,
					address: "来自武汉",
					content: " 仔细一想还真有点感伤[大哭2]",
					like: 3,
					createTime: "1天前"
				},
				{
					id: 16,
					avatar:
						"https://static.juzicon.com/user/avatar-3cb86a0c-08e7-4305-9ac6-34e0cf4937cc-180320123405-BCV6.jpg?x-oss-process=image/resize,m_fill,w_100,h_100",
					parentId: 3,
					username: "Blizzard",
					level: 4,
					content: '回复 <span style="color: blue;"">@别扰我清梦*ぁ:</span> 看完打了一个哈切。。。会传染。。。[委屈]',
					address: "来自广州",
					like: 9,
					createTime: "7天前"
				}
			]
		}
	}
];
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
const tags = ref<Array<string>>(["Tag 1", "Tag 4"]);

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
const scrollHeight = ref("100px");
nextTick(() => {
	fixHeight();
});

const fixHeight = () => {
	scrollHeight.value = document.getElementById("comment-box")!.offsetHeight - 45 + "px";
	console.log("scrollHeight.value", scrollHeight.value);
};
</script>
<style scoped lang="scss">
@import "../index.scss";
</style>
