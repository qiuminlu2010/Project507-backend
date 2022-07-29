<template>
	<!-- 聊天编辑区域 -->
	<div class="rounded-8px">
		<!-- <el-scrollbar height="auto" class="w-100% rounded-8px"> -->
		<!-- <div ref="chatEditor" id="editor-container" @keydown.enter="keyDown()" />
				<input ref="upload" class="opacity-0" @change="handleChange" type="file" /> -->
		<Editor
			style="height: 100px; overflow-y: hidden"
			v-model="valueHtml"
			:defaultConfig="editorConfig"
			mode="simple"
			@keypress.enter="send"
			@onCreated="handleCreated"
			@onChange="handleChange"
		/>
		<div class="float-right my-10px mr-20px">
			<el-button class="w-100px" type="primary" round size="large" color="#e9e9e9" @click="send">发 送</el-button>
		</div>
	</div>
	<!-- </el-scrollbar> -->
</template>
<script setup lang="ts">
// import { ElMessage } from "element-plus";
// import Conversition from "@/class/Conversition";
// import { getCurrentInstance, watch, reactive, onMounted, nextTick } from "vue";
// import { MsgStore } from "..";
// const store = MsgStore();
import { onBeforeUnmount, ref, shallowRef } from "vue";
import { Editor } from "@wangeditor/editor-for-vue";
import "@wangeditor/editor/dist/css/style.css"; // 引入 css
import { Message } from "@/api/interface/index";
import { GlobalStore } from "@/store";
import { MsgStore } from "..";
import { IDomEditor, IEditorConfig } from "@wangeditor/core";
import { Boot } from "@wangeditor/editor";
import ctrlEnterModule from "@wangeditor/plugin-ctrl-enter";
Boot.registerModule(ctrlEnterModule);
// interface EditorProps {
// 	value: string;
// 	id: string;
// 	height: number;
// }
// withDefaults(defineProps<EditorProps>(), {
// 	height: 135,
// 	id: "chatEditor",
// 	value: "value"
// });
// const { proxy }: any = getCurrentInstance();
const globalStore = GlobalStore();
const store = MsgStore();
const editorRef = shallowRef();
const valueHtml = ref("");
const editorConfig: Partial<IEditorConfig> = {
	placeholder: "请输入内容...",
	maxLength: 200
};
// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
	// const editor = editorRef.value;
	const editor = store.editor;
	if (editor == null) return;
	editor.destroy();
});

const handleCreated = (editor: IDomEditor) => {
	// editorRef.value = editor; // 记录 editor 实例，重要！
	store.editor = editor;
	store.editorData = editorRef;
	editorRef.value = editor;
};
const handleChange = () => {
	console.log("change:", store.editor.getHtml().trim());
};

function send() {
	let txt = store.editor.getText();
	if (txt.length !== 0) {
		let msg: Message.MessageInfo = {
			from_uid: globalStore.uid,
			to_uid: store.sessionSelectId,
			content: valueHtml.value,
			image_url: "",
			created_on: Math.round(Date.now() / 1000),
			status: 0
		};
		let wsMsg = {
			content: valueHtml.value,
			to_uid: store.sessionSelectId
		};
		console.log("发送消息", JSON.stringify(wsMsg));
		store.sessionSelected.messages.unshift(msg);
		// store.messageList.get(store.sessionSelectId)?.push(msg);
		// store.toBottom();

		if (store.socket != null) {
			store.socket.send(JSON.stringify(wsMsg));
		}
		setTimeout(() => {
			// store.chatScrollbar.refresh();
			store.chatScrollbar.refresh();
			store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
		}, 1);
		store.editor.setHtml("<p></p>");
	}
}
// 获取文本
// function text() {
// try {
// 	return store.editor.txt.text();
// } catch (e) {
// 	console.warn(e);
// }
// }

// 获取html 有值时为设置内容
// function html(value: any) {
// 	try {
// 		return store.editor.txt.html(value);
// 	} catch (e) {
// 		console.warn(e);
// 	}
// }

// function browsertype() {
// 	let userAgent = navigator.userAgent; //取得浏览器的userAgent字符串
// 	let isOpera = false;
// 	if (userAgent.indexOf("Edge") > -1) {
// 		return "Edge";
// 	}
// 	if (userAgent.indexOf(".NET") > -1) {
// 		return "IE";
// 	}
// 	if (userAgent.indexOf("Opera") > -1 || userAgent.indexOf("OPR") > -1) {
// 		isOpera = true;
// 		return "Opera";
// 	} //判断是否Opera浏览器
// 	if (userAgent.indexOf("Firefox") > -1) {
// 		return "FF";
// 	} //判断是否Firefox浏览器
// 	if (userAgent.indexOf("Chrome") > -1) {
// 		return "Chrome";
// 	}
// 	if (userAgent.indexOf("Safari") > -1) {
// 		return "Safari";
// 	} //判断是否Safari浏览器
// 	if (userAgent.indexOf("compatible") > -1 && userAgent.indexOf("MSIE") > -1 && !isOpera) {
// 		return "IE";
// 	} //判断是否IE浏览器
// }

// 按下回车键
// function keyDown() {
// if (event.ctrlKey && event.keyCode === 13) {
// 	let len = store.editor.txt.html().trim().length;
// 	if (browsertype() == "IE" || browsertype() == "Edge") {
// 		if (len == 0) {
// 			store.editor.cmd.do("insertHTML", "<div></div><div></div>");
// 		} else {
// 			store.editor.cmd.do("insertHTML", "<div></div>");
// 		}
// 	} else if (browsertype() == "FF") {
// 		if (len == 0) {
// 			store.editor.cmd.do("insertHTML", "<br/><br/><br/><br/>");
// 		} else {
// 			store.editor.cmd.do("insertHTML", "<br/><br/>");
// 		}
// 	} else {
// 		if (len == 0) {
// 			store.editor.cmd.do("insertHTML", "<div><br/></div><div><br/></div>");
// 		} else {
// 			store.editor.cmd.do("insertHTML", "<div><br/></div>");
// 		}
// 	}
// } else if (event.keyCode === 13) {
// 	event.preventDefault(); // 阻止浏览器默认换行操作
// 	sendVerify();
// 	return false;
// }
// }

// 发送校验
// function sendVerify() {
// let sendContent = store.editor.txt.html().trim();
// if (!sendContent.length) {
// 	return;
// } else {
// 	let timestamp = +new Date() + "";
// 	let conversition = new Conversition(
// 		store.userInfo.id,
// 		store.recipient.id,
// 		sendContent,
// 		0,
// 		0,
// 		timestamp,
// 		"",
// 		false,
// 		store.userInfo.avatar
// 	);
// 	if (store.socket == null) {
// 		ElMessage.warning("socket不存在");
// 		return;
// 	}
// 	store.sendLocal(conversition);
// 	store.sendInfos(conversition);
// }
// clear();
// }

// 清空
// function clear() {
// 	try {
// 		return store.editor.txt.clear();
// 	} catch (e) {
// 		console.warn(e);
// 	}
// }
</script>

<style scoped lang="scss">
.edit-box {
	.ep-scrollbar {
		border: 1px solid var(--ep-border-color);
	}
	:deep(.w-e-toolbar) {
		border: none !important;
	}
	:deep(.w-e-text-container) {
		height: auto !important;
		min-height: 88px;
		background-color: transparent;
		border: none !important;
	}
}
</style>
