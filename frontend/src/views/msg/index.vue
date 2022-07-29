<template>
	<el-row type="flex" align="middle" justify="center" class="w-full h-full box-border border-rounded-33px">
		<div class="lg:w-1200px flex w-full h-96% mx-20px my-10px box-border border-rounded-12px chat-container">
			<!-- 会话区域 -->
			<ChatDomain />
			<!-- 聊天区域 -->
			<ChatContent v-if="store.sessionSelectId != 0" />
		</div>
	</el-row>
</template>

<script setup lang="ts">
import { getSessionsApi, updateUnreadMessageApi } from "@/api/modules/msg";
import { GlobalStore } from "@/store";
import { onMounted, onUnmounted } from "vue";
import { MsgStore } from ".";
import { Message } from "@/api/interface/index";
import ChatContent from "./components/ChatContent.vue";
import ChatDomain from "./components/ChatDomain.vue";
import { ElMessage } from "element-plus";
import router from "@/routers";
import { onBeforeRouteLeave } from "vue-router";
// import { io } from "socket.io-client";
// import WebSocket from "ws";

const store = MsgStore();
const globalStore = GlobalStore();
// const wsUrl: any = import.meta.env.BASE_URL || "/"

onMounted(() => {
	init();
	initSocket();
});

onBeforeRouteLeave(() => {
	console.log("Connection closed.");
	if (store.sessionSelectId > 0) {
		updateUnreadMessageApi({ uid: globalStore.uid, session_uid: store.sessionSelectId });
	}
});
onUnmounted(() => {});
function initSocket() {
	if (store.socket != null) return;
	let wsUrl = "ws://172.31.225.62:8000/api/v1/msg/" + globalStore.uid + "/chat/" + globalStore.token;
	console.log("初始化ws", wsUrl);
	store.socket = new WebSocket(wsUrl);

	store.socket.onopen = function () {
		console.log("Connection open ...");
		// store.onlineStatus = true;
		// ws.send("Hello WebSockets!");
	};

	store.socket.onmessage = function (evt: any) {
		let msg: Message.MessageInfo = JSON.parse(evt.data);
		console.log("Received Message: ", msg, msg.from_uid);
		if (msg.from_uid !== 0) {
			let item: Message.SessionInfo = store.sessionList.filter((x: Message.SessionInfo) => x.id == msg.from_uid)[0];
			item.messages.unshift(msg);
			if (store.sessionSelectId != msg.from_uid) {
				item.unread += 1;
			} else {
				setTimeout(() => {
					store.chatScrollbar.refresh();
					store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
				}, 10);
			}
		}
	};

	store.socket.onclose = function (evt: any) {
		console.log("Connection closed.", evt);
		if (store.sessionSelectId > 0) {
			updateUnreadMessageApi({ uid: globalStore.uid, session_uid: store.sessionSelectId });
		}
		// if (store.onlineStatus) {
		// 	ElMessage.error("登录失效！");
		// 	globalStore.setToken("");
		// 	router.push({ name: "home" });
		// }
	};
}
async function init() {
	if (globalStore.token) {
		// initSocket();
		if (!store.sessionList.length) {
			console.log("获取会话列表");
			getSessionList(store.sessionListOffset);
		}
	} else {
		ElMessage.error("登录失效！");
		globalStore.setToken("");
		router.push({ name: "home" });
	}
}

// 获取会话列表
async function getSessionList(start: number) {
	let params: Message.ReqGetSessionParams = {
		uid: globalStore.uid,
		page_num: start,
		page_size: 10
	};
	const res = await getSessionsApi(params);
	console.log("获取会话列表", res);
	if (res.code == 200) {
		store.sessionList = res.data!.datalist;
		store.allSessionList = res.data!.datalist;
		// store.sessionList.forEach(item => {
		// 	store.messageList.set(item.id, item.messages);
		// });
	} else {
		ElMessage.error("获取会话列表失败！");
	}

	// userList().then((res: any) => {
	//   if (res?.code == 200) {
	//     store.allSessionList = res?.data
	//   } else {
	//     proxy.$message.error(res?.message)
	//   }
	// })
}
// const screenHeight = computed((): string => (isDevice.value ? window.innerHeight + "px" : "calc(100vh - 80px)"));
</script>
<style scoped lang="scss">
.chat-container {
	box-shadow: 0 1px 2px 0 rgb(0 0 0 / 14%), 0 1px 5px 0 rgb(0 0 0 / 12%);
	:deep(.ep-card__body) {
		padding-left: 0 !important;
		@media (max-width: 1024px) {
			padding-right: 0 !important;
		}
	}
}
</style>
