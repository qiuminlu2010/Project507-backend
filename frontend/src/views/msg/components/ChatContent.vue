<template>
	<!-- 聊天区域 -->
	<div class="h-100% content-box flex box-border">
		<!-- 聊天头部 -->
		<ChatHead></ChatHead>
		<!-- 聊天内容 -->
		<div class="scroll-wrapper rounded-6px box-border w-full message-box max-h-100%" ref="scrollbarRef">
			<div class="px-20px pb-20px">
				<div class="pulldown-wrapper justify-center">
					<div v-if="store.beforePullDown">
						<span>加载更多</span>
					</div>
					<div v-else>
						<div v-if="store.isPullingDown">
							<span>加载中...</span>
						</div>
						<div v-else>
							<span>加载完毕</span>
						</div>
					</div>
				</div>
				<template v-for="(item, index) in messages" :key="index">
					<!-- 聊天时间 -->
					<div class="text-12px flex w-full justify-center">
						<div
							class="bg-gray-200 w-auto text-12px py-5px px-8px text-center rounded-5px"
							v-if="(item && index === 0) || (item && item.created_on - messages[index - 1].created_on > 5 * 60)"
						>
							{{ formatTime(item.created_on) }}
						</div>
					</div>
					<el-row type="flex" class="mb-20px" :class="item.from_uid == globalStore.uid && 'flex-row-reverse'">
						<el-avatar
							:size="42"
							:src="item.from_uid == globalStore.uid ? '/minio' + globalStore.avatar : '/minio' + store.sessionAvatar"
						/>
						<div
							v-if="item.image_url.length == 0"
							:class="item.from_uid == globalStore.uid ? 'mr-10px bg-green-100' : 'ml-10px bg-white'"
							class="break-words px-15px rounded-6px text-left py-12px msg"
							v-dompurify-html="item.content"
						></div>
						<div v-else :class="item.from_uid == globalStore.uid ? 'mr-10px' : 'ml-10px'">
							<el-image
								v-if="item.image_url.length != 0"
								class="w-200px ha max-h-200px"
								:src="item.image_url.indexOf('blob') == 0 ? item.image_url : '/minio' + item.image_url"
								:preview-src-list="[item.image_url.indexOf('blob') == 0 ? item.image_url : '/minio' + item.image_url]"
								:initial-index="4"
								fit="cover"
							/>
							<video v-if="item.image_url.length == 0" width="300" height="180" controls>
								<source :src="item.content" type="video/mp4" />
								您的浏览器不支持 HTML5 video 标签。
							</video>
						</div>
						<el-icon
							v-if="false"
							:class="item.from_uid == globalStore.uid ? 'mr-10px' : 'ml-10px'"
							class="is-loading pt-5px"
							size="40px"
						>
							<i i="ep-loading" class="c-gray-300"></i>
						</el-icon>
					</el-row>
				</template>
			</div>
			<div class="custom-vertical-scrollbar" ref="verticalRef">
				<div class="custom-vertical-indicator"></div>
			</div>
		</div>
		<!-- 聊天底部 -->
		<ChatFoot></ChatFoot>
	</div>
</template>

<script setup lang="ts">
import ChatHead from "./ChatHead.vue";
import ChatFoot from "./ChatFoot.vue";
import { GlobalStore } from "@/store";
import { computed, nextTick, onMounted, ref } from "vue";
// import { computed } from "vue";
import { MsgStore } from "..";
import { formatTime } from "../utils";
import { Message } from "@/api/interface";
import BScroll from "@better-scroll/core";
import MouseWheel from "@better-scroll/mouse-wheel";
import PullDown from "@better-scroll/pull-down";
import { ElMessage } from "element-plus";
import { getMessagesApi } from "@/api/modules/msg";
BScroll.use(PullDown);
BScroll.use(MouseWheel);
const store = MsgStore();
const globalStore = GlobalStore();
const scrollbarRef = ref();
// const verticalRef = ref();

onMounted(() => {
	nextTick(() => {
		// store.chatScrollbar = scrollbarRef.value;
		store.chatScrollbar = new BScroll(".scroll-wrapper", {
			probeType: 3,
			click: false,
			scrollX: false,
			scrollY: true,
			bounceTime: 400,
			useTransition: false,
			pullDownRefresh: {
				threshold: 70,
				stop: 56
			},
			mouseWheel: {
				speed: 20,
				invert: false,
				easeTime: 300,
				dampingFactor: 0.4
			}
		});
		store.chatScrollbar.on("refresh", () => {
			// console.log("会话页面切换");
			// store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
		});
		store.chatScrollbar.on("pullingDown", loadMoreMessages);
	});
	// console.log("加载聊天页面");
	// store.toBottom();
});
// 获取会话列表
const messages = computed(() => {
	return store.sessionSelected.messages.slice().reverse();
});

async function loadMoreMessages() {
	store.beforePullDown = false;
	store.isPullingDown = true;
	console.log("加载更多消息", store.beforePullDown);
	let offset = store.sessionSelected.messages.length;
	let limit = 10;
	// let offset = Math.floor(msgLen / limit);
	// offset += msgLen % limit;
	if (offset < limit || store.sessionSelectId == 0) {
		console.log("结束加载", offset, limit);
		finishPullDown();
		// store.isPullingDown = false;
		// store.beforePullDown = true;
		// store.chatScrollbar.finishPullDown();
		return;
	}
	let _id = store.sessionSelectId;
	let params: Message.ReqGetParams = {
		from_uid: globalStore.uid,
		to_uid: store.sessionSelectId,
		offset: offset,
		limit: limit
	};
	const res = await getMessagesApi(params);
	if (res.code == 200) {
		if (_id === store.sessionSelectId) {
			store.sessionSelected.messages.push(...res.data!.datalist);
		}
		// store.messageList.set(store.sessionSelectId, res.data!.datalist);
		console.log("获取会话消息", "sessionSelectId", store.sessionSelectId, store.sessionSelected.messages);
		// store.initEditor();
	} else {
		ElMessage.error("获取会话消息失败！");
	}
	finishPullDown();
}

async function finishPullDown() {
	store.isPullingDown = false;
	console.log("结束加载");
	// const stopTime = TIME_STOP;
	await new Promise<void>(resolve => {
		setTimeout(() => {
			store.chatScrollbar.finishPullDown();
			resolve();
		}, 600);
	});
	setTimeout(() => {
		store.beforePullDown = true;
		store.chatScrollbar.refresh();
	}, 800);
}
// const messageList = store.messageList.get(store.sessionSelectId);

// 渲染时间每隔5分钟显示一次
// const renderMessageDate = computed(() => {
// 	return (message: any, index: number, messages: Array<any>) => {
// 		console.log("消息时间", message.created_on);
// 		if ((message && index === 0) || (message && message.created_on - messages[index - 1].created_on > 5 * 60 * 1000)) {
// 			return `- - ${formatTime(message.created_on)} - -`;
// 		}
// 		return "";
// 	};
// });
</script>

<style scoped lang="scss">
.content-box {
	position: relative;
	display: flex;
	height: 100%;
	overflow: hidden;
	flex: 1;
	flex-flow: column;
}
.message-box {
	// overflow-y: auto;
	overflow: hidden;

	// height: 100%;
	background: #f8f9fa;
	flex: 1;
}
.emo-image {
	display: inline-block;
	width: 30px;
	height: 30px;
	vertical-align: middle;
}
.msg {
	max-width: calc(70%);
}
.pulldown-wrapper {
	display: flex;

	// position: absolute;
	width: 100%;

	// padding-bottom: 10px;
	color: #999999;
	text-align: center;
	transform: translateY(-100%) translateZ(0);
	box-sizing: border-box;
}

// .custom-vertical-scrollbar {
// 	position: absolute;
// 	top: 50%;
// 	right: 10px;
// 	width: 7px;
// 	height: 200px;
// 	background-color: rgb(200 200 200 / 30%);
// 	border-radius: 6px;
// 	transform: translateY(-20%) translateZ(0);
// }
// .custom-vertical-indicator {
// 	width: 100%;
// 	height: 40px;
// 	background-color: #db8090;
// 	border-radius: 6px;
// }
</style>
