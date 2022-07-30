<template>
	<el-dialog v-model="dialogFormVisible" :title="`添加动态`" :destroy-on-close="true" width="580px" draggable>
		<el-form ref="formRef" :model="form" :rules="rules" class="drawer-multiColumn-form" label-width="100px">
			<el-form-item label="标题 :" prop="title">
				<el-input v-model="form.title" />
			</el-form-item>
			<el-form-item label="内容 :" prop="content">
				<el-input v-model="form.content" :autosize="{ minRows: 2, maxRows: 4 }" type="textarea" />
			</el-form-item>
			<el-form-item label="标签">
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
				<el-button v-else class="btn-tag-input" size="small" @click="showTagInput"> + New Tag </el-button>
			</el-form-item>
			<el-form-item label="数据覆盖 :">
				<el-switch> </el-switch>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" @click="submitForm(formRef)">Create</el-button>
				<el-button @click="resetForm(formRef)">Reset</el-button>
			</el-form-item>
		</el-form>
	</el-dialog>
</template>
<script setup lang="ts">
import { nextTick, reactive, ref } from "vue";
import type { ElInput, FormInstance, FormRules } from "element-plus";

const dialogFormVisible = ref<boolean>(false);

const acceptParams = () => {
	dialogFormVisible.value = true;
};

const rules = reactive<FormRules>({
	title: [
		{ required: true, message: "请输入标题", trigger: "blur" },
		{ min: 3, max: 20, message: "长度必须大于3小于20", trigger: "blur" }
	],
	content: [
		{ required: true, message: "请输入内容", trigger: "blur" },
		{ min: 3, max: 200, message: "长度必须大于3小于200", trigger: "blur" }
	]
});

const formRef = ref<FormInstance>();

const form = reactive({
	title: "",
	content: ""
});

const submitForm = async (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	await formEl.validate((valid, fields) => {
		if (valid) {
			console.log("submit!");
		} else {
			console.log("error submit!", fields);
		}
	});
};

const resetForm = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.resetFields();
};

const tagInputValue = ref("");
const dynamicTags = ref<string[]>([]);
const tagInputVisible = ref(false);
const tagInputRef = ref<InstanceType<typeof ElInput>>();

const showTagInput = () => {
	tagInputVisible.value = true;
	nextTick(() => {
		tagInputRef.value!.input!.focus();
	});
};

const handleTagInputConfirm = () => {
	if (tagInputValue.value) {
		dynamicTags.value!.push(tagInputValue.value);
	}
	tagInputVisible.value = false;
	tagInputValue.value = "";
};

const handleTagClose = (tag: string) => {
	dynamicTags.value!.splice(dynamicTags.value!.indexOf(tag), 1);
};

defineExpose({
	acceptParams
});
</script>

<style scoped lang="scss">
@import "./index.scss";
</style>
