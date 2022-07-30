<template>
	<div>
		<el-dialog v-model="dialogFormVisible" :destroy-on-close="true" width="620px" draggable>
			<el-form ref="formRef" class="drawer-multiColumn-form" label-width="10px">
				<!-- <el-form-item label="标题 :" prop="title">
					<el-input v-model="articleForm.title" />
				</el-form-item> -->
				<el-form-item>
					<el-tag
						v-for="tag in dynamicTags"
						:key="tag"
						class="tag"
						closable
						:disable-transitions="false"
						@close="handleTagClose(tag)"
					>
						{{ tag }}
					</el-tag>
					<el-input
						v-if="tagInputVisible"
						ref="tagInputRef"
						v-model="tagInputValue"
						class="tag-input"
						size="small"
						@keyup.enter="handleTagInputConfirm"
						@blur="handleTagInputConfirm"
					/>
					<el-button v-else class="btn-tag-input" size="small" @click="showTagInput"> + 标签 </el-button>
				</el-form-item>
				<el-form-item>
					<div class="rounded-8px w-full">
						<CommentBox ref="commentRef" placeholder="有什么新鲜事分享给大家" content-btn="发表" :showBtn="false" />
						<div class="flex mt-10px w-full">
							<el-popover placement="top" popper-class="chat-icon-popover" trigger="click">
								<template #reference>
									<div class="text-20px cursor-pointer dark:filter-invert-100">
										<img width="24" height="24" class="block" src="@/assets/images/emoji.svg" />
									</div>
								</template>
								<el-scrollbar class="emoji" height="150px">
									<ul class="m0 p0 pr-2px flex flex-wrap">
										<li
											v-for="(item, row_id) in emojiList"
											:key="row_id"
											class="p-5px list-none hover:animate-heart-beat animate-count-animated animate-duration-1s cursor-pointer"
											:title="item"
											@click="selectIcon"
										>
											{{ item }}
										</li>
									</ul>
								</el-scrollbar>
							</el-popover>
							<div class="i-ep-picture h-24px w-24px ml-12px cursor-pointer" @click="changeUploadImage"></div>
							<div class="i-ep-video-camera h-24px w-24px ml-12px cursor-pointer" @click="changeUploadVideo"></div>
							<div class="i-ep-location h-24px w-24px ml-12px"></div>
							<div class="flex-1"></div>
							<div>
								<el-button class="w-80px" type="primary" round size="large" :disabled="disableBtn" @click="submitForm(formRef)"
									>发 送</el-button
								>
							</div>
						</div>
					</div>
				</el-form-item>
				<el-form-item v-if="isImage"><UploadImg ref="uploadImgRef"></UploadImg></el-form-item>
				<el-form-item v-else>
					<el-upload
						ref="uploadVideoRef"
						class="w-full"
						drag
						action="#"
						multiple
						:auto-upload="false"
						:limit="1"
						:on-exceed="handleExceedVideo"
						:on-change="handleChangeVideo"
					>
						<el-icon class="el-icon--upload"><upload-filled /></el-icon>
						<div class="el-upload__text">拖曳或点击<em>上传视频</em></div>
						<template #tip>
							<div class="el-upload__tip">视频文件要小于 500MB</div>
						</template>
					</el-upload></el-form-item
				>
			</el-form>
		</el-dialog>
	</div>
</template>
<script setup lang="ts">
import { computed, nextTick, ref } from "vue";
import { ElInput, ElMessage, FormInstance, genFileId, UploadFile, UploadProps, UploadRawFile } from "element-plus";
import UploadImg from "@/components/UploadArticle/UploadImg.vue";
// import { Article } from "@/api/interface";
import { GlobalStore } from "@/store";
import { addArticleApi } from "@/api/modules/article";
import CommentBox from "@/components/Comment/CommentBox.vue";
// const imageList = ref();
const tagInputValue = ref("");
const dynamicTags = ref<string[]>([]);
const tagInputVisible = ref(false);
const tagInputRef = ref<InstanceType<typeof ElInput>>();
const dialogFormVisible = ref<boolean>(false);
const commentRef = ref();
const isImage = ref(true);
const formRef = ref<FormInstance>();
const uploadVideoRef = ref();
const uploadImgRef = ref();
let videoFile = ref<UploadFile>();
const disableBtn = computed(() => {
	if (!commentRef.value) return true;
	return commentRef.value.disabled;
});
const acceptParams = () => {
	dialogFormVisible.value = true;
};
// const handleChangePicture: UploadProps["onChange"] = (_uploadFile, uploadFiles) => {
// 	imageList.value = uploadFiles;
// };
// const rules = reactive<FormRules>({
// title: [
// 	{ required: true, message: "请输入标题", trigger: "blur" },
// 	{ min: 3, max: 20, message: "长度必须大于3小于20", trigger: "blur" }
// ],
// content: [
// 	{ required: true, message: "请输入内容", trigger: "blur" },
// 	{ min: 3, max: 200, message: "长度必须大于3小于200", trigger: "blur" }
// ]
// });

// const articleForm = reactive<Article.ReqAddArticleParams>({
// 	title: "",
// 	content: "",
// 	user_id: GlobalStore().uid,
// 	tags: []
// });
// provide("articleForm", articleForm);

const submitForm = async (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	// await formEl.validate(async (valid, fields) => {
	// if (valid) {
	console.log("submit!");
	console.log("content", commentRef.value.content);
	// return;
	let tags = dynamicTags.value;
	// await uploadImgRef.value?.uploadRef.submit();

	let formData = new FormData();

	if (isImage.value && uploadImgRef.value?.files) {
		uploadImgRef.value?.files.forEach((file: UploadFile) => {
			formData.append("images", file.raw!);
		});
	}

	if (!isImage.value && videoFile.value) {
		formData.append("video", videoFile.value.raw!);
		formData.append("type", "1");
	}
	// formData.append("images", param.files);
	tags.forEach(tag => {
		formData.append("tag_name", tag);
	});
	// formData.append("tag_name", articleForm.tags as unknown as Blob);
	// formData.append("title", articleForm.title);
	formData.append("content", commentRef.value.content);
	formData.set("user_id", GlobalStore().uid + "");
	const res = await addArticleApi(formData);
	if (res.code == 200) {
		ElMessage.success("上传成功！");
	} else {
		ElMessage.error("上传失败！");
	}

	console.log(res);
	// } else {
	// console.log("error submit!", fields);
	// return false;
	// }
	dialogFormVisible.value = false;
	// });
};

// const resetForm = (formEl: FormInstance | undefined) => {
// 	if (!formEl) return;
// 	formEl.resetFields();
// };
const showTagInput = () => {
	tagInputVisible.value = true;
	nextTick(() => {
		tagInputRef.value!.input!.focus();
	});
};

const handleTagInputConfirm = () => {
	if (tagInputValue.value) {
		if (dynamicTags.value.indexOf(tagInputValue.value) !== -1) {
			ElMessage.error("标签重复");
		} else {
			dynamicTags.value!.push(tagInputValue.value);
		}
	}
	tagInputVisible.value = false;
	tagInputValue.value = "";
};
function changeUploadImage() {
	isImage.value = true;
	if (!uploadVideoRef.value) return;
	uploadVideoRef.value!.clearFiles();
}
function changeUploadVideo() {
	isImage.value = false;
	if (!uploadImgRef.value) return;
	uploadImgRef.value!.uploadRef.clearFiles();
}
const handleExceedVideo: UploadProps["onExceed"] = files => {
	uploadVideoRef.value!.clearFiles();
	const file = files[0] as UploadRawFile;
	file.uid = genFileId();
	uploadVideoRef.value!.handleStart(file);
};
const handleTagClose = (tag: string) => {
	dynamicTags.value!.splice(dynamicTags.value!.indexOf(tag), 1);
};

const handleChangeVideo: UploadProps["onChange"] = (_uploadFile, uploadFiles) => {
	videoFile.value = uploadFiles[0];
};
function selectIcon() {}
const emojiList =
	"😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊 😫 😴 😌 😛 😜 😒 😓 😔 🙁 😖 😑 😶 🙄 😏 😣 😞 😟 😤 😢 😭 💘 👍 👎 👊 👌".split(" ");
defineExpose({
	acceptParams
});
</script>

<style scoped lang="scss">
@import "./index.scss";
.test {
	background-image: url("@/assets/images/avatar.gif");
}
</style>
