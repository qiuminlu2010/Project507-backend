<template>
	<!-- <ul id="images">
		<li v-for="imgurl in srcList" :key="imgurl"><img :src="imgurl" /></li>
	</ul> -->

	<!-- <lightgallery :settings="gallerySettings" :onInit="onInitGallery" :LG="lg"></lightgallery> -->
	<div id="inline-gallery-container" class="inline-gallery-container"></div>
	<!-- <el-carousel arrow="hover" :autoplay="false" :height="dialogHeight">
		<el-carousel-item v-for="imgurl in srcList" :key="imgurl">
			<el-image style="width: 100%" :src="imgurl" fit="cover" />
		</el-carousel-item>
	</el-carousel> -->
	<!-- <img style="width: 100%" :src="url" /> -->
</template>

<script setup lang="ts" name="cardDialog">
// import { onMounted, ref } from "vue";
// import Lightgallery from "lightgallery/vue";
import lightGallery from "lightgallery";
import "lightgallery/scss/lightgallery.scss";
// import { onMounted } from "vue";
import lgZoom from "lightgallery/plugins/zoom";
import { nextTick } from "vue";
import { CommentStore } from "@/store/modules/comment";
// import Viewer from "viewerjs";
// const imgRef = ref();
const commentStore = CommentStore();
// interface DialogProps {
// 	url: string;
// 	srcList: string[];
// }
// const props = withDefaults(defineProps<DialogProps>(), {
// 	// title: "",
// 	url: "",
// 	srcList: () => []
// });
// const urls = computed(() => {
// 	return commentStore.selectImageUrls;
// });
nextTick(() => {
	// leurls = commentStore.selectImageUrls;
	if (commentStore.gallery === null) {
		let gallery = document.getElementById("inline-gallery-container");
		commentStore.gallery = lightGallery(gallery!, {
			dynamic: true,
			// thumbnail: true,
			container: gallery!,
			plugins: [lgZoom],
			closable: false,
			showMaximizeIcon: true,
			licenseKey: "123",
			dynamicEl: commentStore.selectImageUrls
		});
		commentStore.gallery.openGallery();
	}
});

// const gallerySettings = {
// 	speed: 500,
// 	plugins: [lgZoom],
// 	dynamics: true,
// 	// closable: false,
// 	hash: false,
// 	showMaximizeIcon: true,
// 	dynamicEl: [
// 		{
// 			src: props.srcList[0]
// 		}
// 	]
// };
// let Gallery = lightGallery();
// const onInitGallery = (detail: any) => {
// 	console.log("init gallery");
// };
// Lightgallery.openGallery();
// const dialogHeight = ref("100px");
// onMounted(() => {
// 	fixHeight();
// 	window.onresize = () => {
// 		fixHeight();
// 	};
// });
// const fixHeight = () => {
// 	dialogHeight.value = document.documentElement.clientHeight * 0.96 + "px";
// };
</script>
<style scoped lang="scss">
@import "../index.scss";
.inline-gallery-container {
	width: 100%;

	// set 60% height
	// height: 0;
	// padding-bottom: 65%;
}
</style>
