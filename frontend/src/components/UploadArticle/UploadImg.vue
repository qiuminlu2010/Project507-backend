<template>
	<el-upload
		action="#"
		ref="uploadRef"
		list-type="picture-card"
		:auto-upload="false"
		:on-preview="handlePictureCardPreview"
		:on-remove="handleRemove"
		:limit="excelLimit"
		:before-upload="beforeUpload"
		:on-exceed="handleExceed"
		:on-change="handleChange"
		:http-request="upload"
		accept="image/png,image/jpg,image/jpeg,image/gif"
	>
		<el-icon><Plus /></el-icon>
	</el-upload>

	<el-dialog v-model="dialogVisible">
		<img class="upload-picture" :src="dialogImageUrl" alt="Preview Image" />
	</el-dialog>
</template>

<script lang="ts" setup>
import { ElNotification, UploadFiles, UploadProps } from "element-plus";
import { ref } from "vue";
// import { addArticleApi } from "@/api/modules/article";

const dialogImageUrl = ref("");
const dialogVisible = ref(false);

const handleRemove: UploadProps["onRemove"] = (uploadFile, uploadFiles) => {
	console.log(uploadFile, uploadFiles);
};
const handlePictureCardPreview: UploadProps["onPreview"] = uploadFile => {
	dialogImageUrl.value = uploadFile.url!;
	dialogVisible.value = true;
};
const excelLimit = ref(2);

// 文件上传之前判断
const beforeUpload = (file: any) => {
	const isImg = file.type === "image/png" || file.type === "image/jpg" || file.type === "image/jpeg" || file.type === "image/gif";
	const isLt5M = file.size / 1024 / 1024 < 5;
	if (!isImg)
		ElNotification({
			title: "温馨提示",
			message: "上传文件只能是 png / jpg / jpeg / gif 格式！",
			type: "warning"
		});
	if (!isLt5M)
		ElNotification({
			title: "温馨提示",
			message: "上传文件大小不能超过 5MB！",
			type: "warning"
		});
	return isImg && isLt5M;
};

const handleExceed = (): void => {
	ElNotification({
		title: "温馨提示",
		message: "最多只能上传九张图片！",
		type: "warning"
	});
};

let files = ref<UploadFiles>();
const handleChange: UploadProps["onChange"] = (_uploadFile, uploadFiles) => {
	files.value = uploadFiles;
};

const uploadRef = ref();

// const articleForm = inject("articleForm") as Article.ReqAddArticleWithImgParams;
// 文件上传
const upload = async () => {
	// let formData = new FormData();
	// console.log("uploadFiles", files);
	// files.value.forEach(file => {
	// 	formData.append("images", file.raw);
	// });
	// articleForm.tags.forEach(tag => {
	// 	formData.append("tag_name", tag);
	// });
	// formData.append("title", articleForm.title);
	// formData.append("content", articleForm.content);
	// formData.set("user_id", articleForm.user_id);
	// const res = await addArticleApi(formData);
	// console.log(res);
	// dialogVisible.value = false;
};
defineExpose({
	files,
	uploadRef
});
</script>

<style scoped lang="scss">
@import "./index.scss";
</style>
