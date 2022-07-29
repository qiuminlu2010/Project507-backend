<template>
	<div class="table-box">
		<ProTable ref="proTable" :requestApi="getArticleListApi" :initParam="initParam" :columns="columns">
			<!-- 表格 header 按钮 -->
			<template #tableHeader="scope">
				<el-button type="primary" :icon="CirclePlus" @click="addNew">新增文章</el-button>
				<el-button type="primary" :icon="Upload" plain @click="null">批量添加文章</el-button>
				<el-button type="primary" :icon="Download" plain @click="downloadFile">导出文章数据</el-button>
				<el-button type="danger" :icon="Delete" plain :disabled="!scope.isSelected" @click="null"> 批量删除用户 </el-button>
			</template>
			<!-- Expand -->
			<template #expand="scope">
				{{ scope.row }}
			</template>
			<!-- 文章状态 slot -->
			<template #state="scope">
				<!-- 如果插槽的值为 el-switch，第一次加载会默认触发 switch 的 @change 方法，所有在外层包一个盒子，点击触发盒子 click 方法 -->
				<div @click="changeStatus(scope.row)">
					<el-switch
						:value="scope.row.state"
						:active-text="scope.row.state === 1 ? '启用' : '禁用'"
						:active-value="1"
						:inactive-value="0"
					/>
				</div>
			</template>
			<!-- 表格操作 -->
			<template #operation="scope">
				<el-button type="primary" link :icon="View" @click="openDrawer('查看', scope.row)">查看</el-button>
				<el-button type="primary" link :icon="EditPen" @click="openDrawer('编辑', scope.row)">编辑</el-button>
				<el-button type="primary" link :icon="Delete" @click="deleteItem(scope.row)">删除</el-button>
			</template>
		</ProTable>
		<UserDrawer ref="drawerRef"></UserDrawer>
		<UploadArticle ref="dialogRef"></UploadArticle>
	</div>
</template>

<script setup lang="tsx" name="articleList">
import { ref, reactive, provide } from "vue";
// import { genderType } from "@/utils/serviceDict";
import { ElMessage } from "element-plus";
import { Article, User } from "@/api/interface";
import { ColumnProps } from "@/components/ProTable/interface";
import { useHandleData } from "@/hooks/useHandleData";
import { useDownload } from "@/hooks/useDownload";
import ProTable from "@/components/ProTable/index.vue";
import UploadArticle from "@/components/UploadArticle/index.vue";
import UserDrawer from "@/views/admin/components/UserDrawer.vue";
import { CirclePlus, Delete, EditPen, Download, Upload, View } from "@element-plus/icons-vue";
import {
	getArticleListApi,
	changeArticleStatus,
	deleteArticle,
	exportArticleInfo,
	addArticleApi,
	editArticle
} from "@/api/modules/article";

// 获取 ProTable 元素，调用其获取刷新数据方法（还能获取到当前查询参数，方便导出携带参数）
const proTable = ref();
// 如果表格需要初始化请求参数,直接定义传给 ProTable(之后每次请求都会自动带上该参数，此参数更改之后也会一直带上)
const initParam = reactive({
	type: 1
});

// 自定义渲染头部(使用tsx语法)
const renderHeader = (scope: any) => {
	return (
		<el-button
			type="primary"
			onClick={() => {
				console.log(scope);
				ElMessage.success("我是自定义表头");
			}}
		>
			{scope.row.label}
		</el-button>
	);
};

// 配置项
const columns: Partial<ColumnProps>[] = [
	{
		type: "selection",
		width: 80,
		fixed: "left"
	},
	{
		type: "index",
		label: "#",
		width: 80
	},
	{
		prop: "ID",
		label: "ID",
		sortable: true,
		search: true
	},
	{
		prop: "title",
		label: "标题",
		sortable: true,
		search: true,
		width: 160,
		renderHeader
	},
	{
		prop: "content",
		label: "内容",
		search: true,
		width: 160
	},
	{
		prop: "tag",
		label: "标签s",
		width: 100
	},
	{
		prop: "user_id",
		label: "用户ID",
		search: true,
		sortable: true
	},
	{
		prop: "watch",
		label: "点击量",
		search: true,
		sortable: true
	},
	{
		prop: "like",
		label: "点赞",
		search: true,
		sortable: true
	},
	{
		prop: "collect",
		label: "收藏",
		search: true,
		sortable: true
	},
	{
		prop: "state",
		label: "状态",
		width: 120
	},
	{
		prop: "thumb",
		label: "封面",
		width: 120,
		image: true,
		renderHeader
	},
	{
		prop: "created_on",
		label: "创建时间",
		width: 200,
		sortable: true,
		search: true,
		searchType: "datetimerange",
		initSearchParam: ["2022-04-05 00:00:00", "2022-05-10 23:59:59"]
	},
	{
		prop: "modified_on",
		label: "修改时间",
		width: 200,
		sortable: true,
		search: true,
		searchType: "datetimerange",
		initSearchParam: ["2022-04-05 00:00:00", "2022-05-10 23:59:59"]
	},
	{
		prop: "operation",
		label: "操作",
		width: 320,
		fixed: "right"
	}
];

const deleteItem = async (row: Article.ResArticleList) => {
	await useHandleData(deleteArticle, { id: row.id }, `删除文章【${row.title}】`);
	proTable.value.refresh();
};

// 导出用户列表
const downloadFile = async () => {
	useDownload(exportArticleInfo, "文章列表", proTable.value.searchParam);
};
// 切换状态
const changeStatus = async (row: Article.ResArticleList) => {
	await useHandleData(changeArticleStatus, { id: row.id, state: row.state == 1 ? 0 : 1 }, `切换【${row.title}】用户状态`);
	proTable.value.refresh();
};

interface DialogExpose {
	acceptParams: () => void;
}
const dialogRef = ref<DialogExpose>();
const addNew = () => {
	dialogRef.value!.acceptParams();
};
// 打开 drawer(新增、查看、编辑)
interface DrawerExpose {
	acceptParams: (params: any) => void;
}
const drawerRef = ref<DrawerExpose>();
const openDrawer = (title: string, rowData: Partial<User.ResUserList> = {}) => {
	let params = {
		title: title,
		rowData: { ...rowData },
		isView: title === "查看" ? true : false,
		apiUrl: title === "新增" ? addArticleApi : title === "编辑" ? editArticle : "",
		getTableList: proTable.value.refresh
	};
	drawerRef.value!.acceptParams(params);
};
provide("proTable", proTable);
</script>
