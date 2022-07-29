<template>
	<el-drawer v-model="drawerVisible" :destroy-on-close="true" size="600px" :title="`${drawerData.title}用户`">
		<el-form ref="ruleFormRef" :rules="rules" :disabled="drawerData.isView" :model="drawerData.rowData" label-width="100px">
			<el-form-item label="用户名" prop="username" clearable>
				<el-input v-model="drawerData.rowData!.username" placeholder="请填写用户名"></el-input>
			</el-form-item>
			<el-form-item label="密码" prop="password" clearable>
				<el-input v-model="drawerData.rowData!.password" placeholder="请填写密码"></el-input>
			</el-form-item>
			<el-form-item label="校园卡号" prop="student_id" clearable>
				<el-input v-model="drawerData.rowData!.student_id" placeholder="请填写校园卡号"></el-input>
			</el-form-item>
			<el-form-item label="宿舍楼" prop="address" clearable>
				<el-input v-model="drawerData.rowData!.address" placeholder="请填写宿舍楼"></el-input>
			</el-form-item>
		</el-form>
		<template #footer>
			<el-button @click="drawerVisible = false">取消</el-button>
			<el-button type="primary" v-show="!drawerData.isView" @click="handleSubmit">确定</el-button>
		</template>
	</el-drawer>
</template>

<script setup lang="ts" name="UserDrawer">
import { User } from "@/api/interface";
import { ref, reactive } from "vue";
// import { genderType } from "@/utils/serviceDict";
import { ElMessage, FormInstance } from "element-plus";

const rules = reactive({
	username: [{ required: true, message: "请填写用户名", trigger: "blur" }],
	password: [{ required: true, message: "请填写用户名", trigger: "blur" }],
	student_id: [{ required: false, message: "请填写校园卡号", trigger: "blur" }],
	address: [{ required: false, message: "请填写宿舍楼", trigger: "blur" }]
});

interface DrawerProps {
	title: string;
	isView: boolean;
	rowData?: User.ResUserList;
	apiUrl?: (params: any) => Promise<any>;
	getTableList?: () => Promise<any>;
}

// drawer框状态
const drawerVisible = ref(false);
const drawerData = ref<DrawerProps>({
	isView: false,
	title: ""
});

// 接收父组件传过来的参数
const acceptParams = (params: DrawerProps): void => {
	drawerData.value = params;
	drawerVisible.value = true;
};

const ruleFormRef = ref<FormInstance>();
// 提交数据（新增/编辑）
const handleSubmit = () => {
	ruleFormRef.value!.validate(async valid => {
		if (!valid) return;
		try {
			if (!drawerData.value.apiUrl) return;
			await drawerData.value.apiUrl(drawerData.value.rowData);
			ElMessage.success({ message: `${drawerData.value.title}用户成功！` });
			drawerData.value.getTableList && drawerData.value.getTableList();
			drawerVisible.value = false;
		} catch (error) {
			console.log(error);
		}
	});
};

defineExpose({
	acceptParams
});
</script>
