<template>
	<div>
		<Comment></Comment>
	</div>
</template>

<script setup lang="ts">
// import { reactive } from "vue";
import Comment from "@/components/Comment/index.vue";
// import { CommentApi } from "@/components/Comment/interface";
import { CommentStore } from "@/store/modules/comment";
// 下载表情包资源emoji.zip https://gitee.com/undraw/undraw-ui/releases
// static文件放在public下,引入emoji.ts文件可以移动到自定义位置
// import emoji from "@/assets/emoji";
// const emoji: never[] = [];
// const comments = ref([] as CommentApi[]);
const commentStore = CommentStore();
// const user = reactive({
// 	id: 1,
// 	username: "user",
// 	avatar: "https://static.juzicon.com/avatars/avatar-200602130320-HMR2.jpeg?x-oss-process=image/resize,w_100",
// 	// 评论id数组
// 	likes: [1, 2, 11]
// });

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
</script>
