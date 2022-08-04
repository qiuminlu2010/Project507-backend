<template>
	<!-- <ul id="images">
		<li v-for="imgurl in srcList" :key="imgurl"><img :src="imgurl" /></li>
	</ul> -->

	<!-- <lightgallery :settings="gallerySettings" :onInit="onInitGallery" :LG="lg"></lightgallery> -->
	<div class="w-full justify-center flex"><vue3VideoPlay v-bind="options" :poster="commentStore.selectPreviewUrl" /></div>
	<!-- <el-carousel arrow="hover" :autoplay="false" :height="dialogHeight">
		<el-carousel-item v-for="imgurl in srcList" :key="imgurl">
			<el-image style="width: 100%" :src="imgurl" fit="cover" />
		</el-carousel-item>
	</el-carousel> -->
	<!-- <img style="width: 100%" :src="url" /> -->
</template>

<script setup lang="ts" name="cardDialog">
import { CommentStore } from "@/store/modules/comment";
import { computed, reactive } from "vue";
const commentStore = CommentStore();
const videoSrc = computed(() => {
	return `/minio` + commentStore.selectVideoUrl;
});
const options = reactive({
	width: "100%", //播放器高度
	height: "100%", //播放器高度
	color: "#409eff", //主题色
	title: "", //视频名称
	src: videoSrc, //视频源
	type: "m3u8",
	muted: false, //静音
	webFullScreen: false,
	speedRate: ["0.75", "1.0", "1.25", "1.5", "2.0"], //播放倍速
	autoPlay: false, //自动播放
	loop: false, //循环播放
	mirror: false, //镜像画面
	ligthOff: false, //关灯模式
	volume: 0.3, //默认音量大小
	control: true, //是否显示控制
	controlBtns: ["audioTrack", "quality", "speedRate", "volume", "setting", "pip", "pageFullScreen", "fullScreen"] //显示所有按钮,
});
</script>
<style scoped lang="scss"></style>
