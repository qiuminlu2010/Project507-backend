<template>
	<el-dropdown trigger="click">
		<div class="avatar">
			<!-- <img :src="'/base' + globalStore.avatar" alt="avatar" /> -->
			<avatar :username="globalStore.username" :size="50" :src="'/base' + globalStore.avatar"></avatar>
		</div>
		<template #dropdown>
			<el-dropdown-menu>
				<el-dropdown-item @click="openDialog('infoRef')">{{ $t("header.personalData") }}</el-dropdown-item>
				<el-dropdown-item @click="openDialog('uploadAvatarRef')">修改头像</el-dropdown-item>
				<el-dropdown-item @click="logout" divided>{{ $t("header.logout") }}</el-dropdown-item>
			</el-dropdown-menu>
		</template>
	</el-dropdown>
	<!-- infoDialog -->
	<InfoDialog ref="infoRef"></InfoDialog>
	<!-- passwordDialog -->
	<UploadAvatar ref="uploadAvatarRef"></UploadAvatar>
</template>

<script setup lang="ts">
import { ref } from "vue";
import InfoDialog from "./InfoDialog.vue";
// import PasswordDialog from "./PasswordDialog.vue";
import { ElMessageBox, ElMessage } from "element-plus";
import { useRouter } from "vue-router";
import { GlobalStore } from "@/store";
import UploadAvatar from "@/components/UploadAvatar/index.vue";
const router = useRouter();
const globalStore = GlobalStore();

// logout
const logout = () => {
	ElMessageBox.confirm("您是否确认退出登录?", "温馨提示", {
		confirmButtonText: "确定",
		cancelButtonText: "取消",
		type: "warning"
	}).then(() => {
		router.push({ name: "login" });
		globalStore.setToken("");
		ElMessage({
			type: "success",
			message: "退出登录成功！"
		});
	});
};

interface DialogExpose {
	openDialog: () => void;
}
const infoRef = ref<null | DialogExpose>(null);
const uploadAvatarRef = ref<null | DialogExpose>(null);
// openDialog
const openDialog = (refName: string) => {
	if (refName == "infoRef") return infoRef.value?.openDialog();
	uploadAvatarRef.value?.openDialog();
};
</script>

<style scoped lang="scss">
@import "../index.scss";
</style>
