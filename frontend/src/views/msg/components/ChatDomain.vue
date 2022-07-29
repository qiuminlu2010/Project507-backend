<template>
	<div class="chat-session p-8px box-border">
		<div class="session-search-box">
			<div class="search-icon">
				<Search style="width: 1em; height: 1em; margin-right: 8px" />
			</div>
			<input
				v-model="searchInput"
				type="search"
				placeholder="搜索"
				autocomplete="off"
				class="search-input"
				spellcheck="false"
				data-ms-editor="true"
			/>
		</div>
		<el-scrollbar height="auto">
			<!-- 会话列表 -->
			<template v-for="item in handleSessionList" :key="item">
				<div
					class="session-item cursor-pointer rounded-8px"
					:class="[item.id === store.sessionSelectId ? 'session-active' : '']"
					@click="selectSession(item)"
				>
					<div class="session-container">
						<el-badge class="mr-15px" :value="item.unread" :max="99" :hidden="item.unread == 0">
							<el-avatar shape="square" class="!block" :size="42" fit="cover" :src="item.avatar" />
						</el-badge>
						<div class="w-100% text-left text-container text-ellipsis">
							<div class="truncate h-25px text-title">
								<div class="username text-ellipsis">{{ item.username }}</div>
								<div class="text-date">{{ formatTime(item.messages[0].created_on) }}</div>
							</div>
							<div
								class="truncate text-12px h-16px max-w-258px"
								style="color: #67717a"
								v-dompurify-html="item.messages[0].content"
							></div>
						</div>
					</div>
				</div>
			</template>
			<el-empty v-if="!handleSessionList.length" :image-size="100" description="消息列表为空" />
		</el-scrollbar>
	</div>
</template>

<script setup lang="ts">
import { Message } from "@/api/interface";
import { updateUnreadMessageApi } from "@/api/modules/msg";
import { GlobalStore } from "@/store";
// import { ElMessage } from "element-plus";
import { computed, ref, watch } from "vue";
import { MsgStore } from "..";
import { Search } from "@element-plus/icons-vue";
import { formatTime } from "../utils";
const store = MsgStore();
const globalStore = GlobalStore();
const searchInput = ref("");

watch(
	() => searchInput.value,
	() => {
		console.log("会话列表搜索");
		if (searchInput.value.length == 0) {
			store.sessionList = store.allSessionList.slice();
		} else {
			store.sessionList = store.allSessionList
				.filter((x: Message.SessionInfo) => x.username.indexOf(searchInput.value) != -1)
				.slice();
		}
		// setTimeout(() => {
		// 	// store.chatScrollbar.refresh();
		// 	store.chatScrollbar.refresh();
		// 	store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
		// }, 100);
	}
);
// 选择聊天用户
async function selectSession(item: Message.SessionInfo) {
	// 更新已读消息状态
	// store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
	// store.chatScrollbar.refresh();
	if (store.sessionSelectId > 0 && store.sessionSelected !== null && store.sessionSelected.unread > 0) {
		updateUnreadMessageApi({ uid: globalStore.uid, session_uid: store.sessionSelectId });
	}
	store.sessionSelectId = item.id;
	if (store.sessionSelectId > 0 && item.unread > 0) {
		updateUnreadMessageApi({ uid: globalStore.uid, session_uid: store.sessionSelectId });
	}
	if (item.avatar.length == 0) {
		// store.sessionAvatar = "@/assets/images/avatar.png";
		store.sessionAvatar = "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png";
	} else {
		store.sessionAvatar = item.avatar;
	}
	store.sessionName = item.username;
	store.sessionSelected = item;
	item.unread = 0;
	// store.chatScrollbar.refresh();
	// store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
	setTimeout(() => {
		// store.chatScrollbar.refresh();
		store.chatScrollbar.refresh();
		store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
	}, 10);
	setTimeout(() => {
		// store.chatScrollbar.refresh();
		store.chatScrollbar.refresh();
		store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
	}, 10);
	// if (!store.messageList.has(store.sessionSelectId)) {
	// 	let params: Message.ReqGetParams = {
	// 		from_uid: globalStore.uid,
	// 		to_uid: store.sessionSelectId,
	// 		pageNum: 0,
	// 		pageSize: 10
	// 	};
	// 	const res = await getMessagesApi(params);
	// 	if (res.code == 200) {
	// 		store.messageList.set(store.sessionSelectId, res.data!.datalist);
	// 		console.log("获取会话消息", "sessionSelectId", store.sessionSelectId, store.messageList.get(store.sessionSelectId));
	// 		// store.initEditor();
	// 	} else {
	// 		ElMessage.error("获取会话消息失败！");
	// 	}
	// }
	// store.toBottom();
}
// 获取最后一条消息内容
// const getLastSession = computed(() => {
// 	// return item.messages[0];
// 	return "最后一条信息";
// });
// 返回对应选择列表
const handleSessionList = computed(() => {
	console.log("返回会话列表", store.sessionList);
	return store.sessionList;
});
// 获取未读数量
// const getUnReadCount = null;
</script>

<style scoped lang="scss">
/* stylelint-disable-next-line scss/dollar-variable-pattern */
$selectBg: #f6f6f6;
.chat-session {
	// min-width: 260px;
	// max-width: 500px;
	position: relative;
	display: flex;
	height: 100%;
	background: #ffffff;
	border-right: 1px solid #e1e5e8;
	flex-flow: column;
	flex: 0 0 30%;
	:deep(p) {
		margin: 0;
	}
	:deep(.el-scrollbar) {
		height: auto;
	}

	// @media (max-width: 1024px) {
	// 	border-right: none;
	// }
	.session-search-box {
		position: sticky;
		display: flex;
		align-items: center;
		height: 64px;
		padding: 0 15px;
		.search-icon {
			position: absolute;
			left: 30px;
			display: flex;
			color: #a8aeb3;
			svg {
				width: 18px;
				height: 18px;
			}
		}
		.search-input {
			width: 100%;
			height: 38px;
			min-width: 10px;
			padding: 10px 10px 10px 40px;
			font-size: 15px;
			color: #0a0a0a;
			background: #ffffff;
			border: 1px solid #e1e5e8;
			border-radius: 4px;
			border-radius: 20px;
			outline: 0;
			caret-color: #0a0a0a;
		}
	}
	.session-item {
		position: relative;
		display: flex;
		align-items: center;
		min-height: 71px;
		padding: 0 14px;
		margin-bottom: 5px;
		transition: background-color 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
		flex: 1 1 100%;
		&:hover {
			background: $selectBg;
		}
		.text-ellipsis {
			width: 100%;
			overflow: hidden;
			text-overflow: ellipsis;
			white-space: nowrap;
		}
		.session-container {
			display: flex;
			flex: 1;
			align-items: center;
			width: 100%;
			.text-container {
				flex: 1;
			}
			.text-title {
				display: flex;
				align-items: center;
				line-height: 25px;
				.username {
					font-weight: 500;
					color: #0a0a0a;
					flex: 1;
				}
				.text-date {
					margin-left: 5px;
					font-size: 11px;
					color: #a8aeb8;
				}
			}
		}
	}
	.session-active {
		background: #e5effa;
	}
}
</style>
