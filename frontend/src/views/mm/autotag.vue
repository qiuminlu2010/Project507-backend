<template>
	<div class="main">
		<div class="upload-demo flex flex-col">
			<div class="title">AutoTag</div>
			<div class="flex justify-center">
				<el-upload ref="uploadRef" action="#" :http-request="uploadAction" :show-file-list="false">
					<template #trigger>
						<el-button type="primary">上传图片</el-button>
					</template>
					<!-- <template #tip>
						<div class="el-upload__tip">支持jpg/png 格式 文件小于5MB</div>
					</template> -->
				</el-upload>
			</div>
			<el-image class="display" v-if="imageUrl" :src="imageUrl" fit="contain" />
			<div class="p-2% pb-15%">
				<el-tag
					v-for="tag in tags"
					:key="tag"
					class="mx-1"
					size="large"
					closable
					:disable-transitions="false"
					@close="handleClose(tag)"
					round
				>
					{{ tag }}
				</el-tag>
			</div>
		</div>
		<!-- <img class="display_img" v-if="imageUrl" :src="imageUrl" /> -->
	</div>
</template>
<script setup lang="ts">
// import { UploadFilled } from "@element-plus/icons-vue";
import axios from "axios";
// import type { UploadProps } from "element-plus";
import { ref } from "vue";
const imageUrl = ref("");
const uploadRef = ref();
// const handleAvatarSuccess: UploadProps["onSuccess"] = response => {
// 	imageUrl.value = URL.createObjectURL(response);
// };
const uploadAction = (option: any) => {
	let param = new FormData();
	param.append("file", option.file);
	// let config = {
	// 	headers: { "Content-Type": "multipart/form-data" }
	// };
	imageUrl.value = URL.createObjectURL(option.file);
	axios.post("/flare/autotag", param).then(response => {
		console.log(response);
		tags.value = response.data["tags"];
		// imageUrl.value = URL.createObjectURL(response.data);
		// this.imgUrl = window.URL.createObjectURL(res.data); //这里调用window的URL方法
	});
};

const tags = ref<string[]>([]);

const handleClose = (tag: string) => {
	tags.value.splice(tags.value.indexOf(tag), 1);
};
</script>
<style scoped lang="scss">
:deep(.el-upload-dragger) {
	background-color: rgba($color: #e3e5e7, $alpha: 10%);
}
.main {
	display: flex;
	justify-content: center;
	width: 100%;
	height: 100%;
	// background: #000000;

	// justify-items: center;
}
.title {
	padding: 5% 3%;
	font-size: 44px;
	font-weight: bold;
	// color: #ffffff;
	text-align: center;
}
.upload-demo {
	width: 50%;
	height: 95%;
	margin: 1%;
	border: 1px solid #000000;
	border-radius: 1rem;
	box-sizing: border-box;

	// padding: 1%;
	// margin-top: 10%;

	// margin-left: 25%;
}
.upload_box {
	// padding: 5%;
	// margin-top: 10%;
}
.display {
	// width: 70%;
	// height: 98%;
	padding: 2%;
	box-sizing: border-box;

	// .display_img {
	// 	// object-fit: contain;
	// }
}
</style>
