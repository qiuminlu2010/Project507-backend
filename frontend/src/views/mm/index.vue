<template>
	<div>
		<vue3VideoPlay v-bind="options" poster="https://cdn.jsdelivr.net/gh/xdlumia/files/video-play/ironMan.jpg" />
	</div>
</template>
<script setup lang="ts">
import { reactive } from "vue";
const options = reactive({
	width: "800px", //播放器高度
	height: "450px", //播放器高度
	color: "#409eff", //主题色
	title: "", //视频名称
	src: "http://192.168.198.132:32000/video/1c938d0a937d590863d5fe453ae7db91.m3u8", //视频源
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

<style scoped></style>
