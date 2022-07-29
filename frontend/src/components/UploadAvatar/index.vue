<template>
	<my-upload
		field="img"
		v-model="dialogVisible"
		@crop-success="cropSuccess"
		@crop-upload-success="cropUploadSuccess"
		@crop-upload-fail="cropUploadFail"
		:width="200"
		:height="200"
		method="PUT"
		:headers="{ token: store.token }"
		:url="'/base/api/v1/user/' + store.uid + '/avatar'"
		img-format="png"
	></my-upload>
</template>

<script setup lang="ts">
import { GlobalStore } from "@/store";
import { ElMessage } from "element-plus";
import { ref } from "vue";
import myUpload from "vue-image-crop-upload/upload-3.vue";
// import "babel-polyfill";
const dialogVisible = ref(false);
const store = GlobalStore();
// openDialog
const openDialog = () => {
	dialogVisible.value = true;
};

function cropSuccess() {
	console.log("准备提交");
}
/**
 * upload success
 *
 * [param] jsonData  server api return data, already json encode
 * [param] field
 */
function cropUploadSuccess() {
	ElMessage.success("上传成功");
	console.log("上传成功");
}
/**
 * upload fail
 *
 * [param] status    server api return error status, like 500
 * [param] field
 */
function cropUploadFail() {
	ElMessage.error("头像上传失败");
}
defineExpose({
	openDialog
});
</script>
