<template>
	<div class="session-footer border-box">
		<div class="pl-10px mt-20px">
			<!--工具栏-->
			<el-row type="flex" class="mb-10px pl-10px">
				<el-popover placement="top" popper-class="chat-icon-popover" trigger="click">
					<template #reference>
						<div class="text-20px cursor-pointer dark:filter-invert-100">
							<img width="24" height="24" class="block" src="@/assets/images/emoji.svg" />
						</div>
					</template>
					<el-scrollbar class="emoji" height="150px">
						<ul class="m0 p0 pr-2px flex flex-wrap">
							<li
								v-for="(item, row_id) in store.emojiList"
								:key="row_id"
								class="p-5px list-none hover:animate-heart-beat animate-count-animated animate-duration-1s cursor-pointer"
								:title="item"
								@click="selectIcon(item)"
							>
								{{ item }}
							</li>
						</ul>
					</el-scrollbar>
				</el-popover>
				<div class="ml-10px text-20px i-ep-picture-rounded !cursor-pointer">
					<input
						ref="UploadImageRef"
						class="opacity-0"
						name="customerService"
						type="file"
						value=""
						accept="image/*"
						@change="sendImage"
					/>
				</div>
				<div class="ml-10px text-20px i-ep-video-camera !cursor-pointer">
					<input
						ref="referenceUploadVideo"
						class="opacity-0 hidden"
						name="customerService"
						type="file"
						value=""
						accept="video/*"
						@change="sendVideo"
					/>
				</div>
			</el-row>
			<ChatEditor ref="editor"></ChatEditor>
		</div>
	</div>
</template>

<script setup lang="ts">
import { MsgStore } from "..";
import ChatEditor from "./ChatEditor.vue";
import { Message } from "@/api/interface/index";
import { ElMessage } from "element-plus";
import { getCurrentInstance, ref } from "vue";
import { GlobalStore } from "@/store";
import { uploadImageApi, uploadVideoApi } from "@/api/modules/msg";
// import { ElMessage } from "element-plus/lib/components";
// const { proxy }: any = getCurrentInstance();
const store = MsgStore();
const globalStore = GlobalStore();
const { proxy }: any = getCurrentInstance();
const editor = ref(null);
// 选择表情
const selectIcon = (item: any) => {
	// const editor = store.editor;
	// if (editor == null) return;
	// console.log("选择表情1", item, editor.getHtml());
	store.insertText(item);
	store.editor.focus();
};

// function blurHighLight() {}
//发送图片
async function sendImage(e: any) {
	const fileData = e.target.files[0];
	if (fileData != null) {
		if (!/image\/\w+/.test(fileData.type)) {
			return alert("请选择图片文件!");
		}
		if (fileData.size > 1024 * 1024 * 10) {
			return alert("上传图片不能超过10M!");
		} else {
			console.log(fileData);
			const tempFilePath = URL.createObjectURL(fileData);
			//localMsg
			let msg: Message.MessageInfo = {
				from_uid: globalStore.uid,
				to_uid: store.sessionSelectId,
				content: "",
				image_url: tempFilePath,
				created_on: Math.round(Date.now() / 1000),
				status: 0
			};
			let formData = new FormData();
			formData.append("image", fileData);
			const res = await uploadImageApi(formData);
			if (res.code === 200) {
				//添加本地消息
				store.sessionSelected.messages.unshift(msg);
				let wsMsg = {
					content: "",
					image_url: res.data!.image_url,
					to_uid: store.sessionSelectId
				};
				//发送给对方消息
				if (store.socket != null) {
					store.socket.send(JSON.stringify(wsMsg));
				}
			} else {
				ElMessage.error("发送图片失败！");
			}
			setTimeout(() => {
				// store.chatScrollbar.refresh();
				store.chatScrollbar.refresh();
				store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
			}, 1);
			proxy.$refs.UploadImageRef.value = null;
		}
	}
}
async function sendVideo(e: any) {
	const resultFile = e.target.files;
	const fileObj = new Blob([resultFile[0]], { type: "video/mp4" });
	const tempFilePath = URL.createObjectURL(fileObj);
	let msg: Message.MessageInfo = {
		from_uid: globalStore.uid,
		to_uid: store.sessionSelectId,
		content: "",
		image_url: "",
		video_url: tempFilePath,
		created_on: Math.round(Date.now() / 1000),
		status: 0
	};
	let wsMsg = {
		content: "",
		image_url: "",
		to_uid: store.sessionSelectId
	};
	store.sessionSelected.messages.push(msg);
	let formData = new FormData();
	formData.append("video", resultFile[0]);
	const res = await uploadVideoApi(formData);
	if (res.code === 200) {
		wsMsg.image_url = res.data!.video_url;
	}
	if (store.socket != null) {
		store.socket.send(JSON.stringify(wsMsg));
	}
	setTimeout(() => {
		// store.chatScrollbar.refresh();
		store.chatScrollbar.refresh();
		store.chatScrollbar.scrollTo(0, store.chatScrollbar.maxScrollY);
	}, 1);
	proxy.$refs.UploadImageRef.value = null;
}

store.emojiList =
	"😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊 😫 😴 😌 😛 😜 😒 😓 😔 🙁 😖 😑 😶 🙄 😏 😣 😞 😟 😤 😢 😭 💘 👍 👎 👊 👌".split(" ");
</script>
<style scoped lang="scss">
.session-footer {
	z-index: 10;
	width: 100%;
	border-top: 1px solid #e1e5e8;
}
</style>
