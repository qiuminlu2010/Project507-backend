import { RouteRecordRaw } from "vue-router";
import { LayoutAdmin } from "@/routers/constant";

// 表单 Form 模块
const formRouter: Array<RouteRecordRaw> = [
	{
		path: "/form",
		component: LayoutAdmin,
		redirect: "/form/basicForm",
		meta: {
			title: "表单 Form"
		},
		children: [
			{
				path: "/form/basicForm",
				name: "basicForm",
				component: () => import("@/views/form/basicForm/index.vue"),
				meta: {
					keepAlive: true,
					requiresAuth: true,
					title: "基础 Form",
					key: "basicForm"
				}
			},
			{
				path: "/form/validateForm",
				name: "validateForm",
				component: () => import("@/views/form/validateForm/index.vue"),
				meta: {
					keepAlive: true,
					requiresAuth: true,
					title: "校验 Form",
					key: "validateForm"
				}
			},
			{
				path: "/form/dynamicForm",
				name: "dynamicForm",
				component: () => import("@/views/form/dynamicForm/index.vue"),
				meta: {
					keepAlive: true,
					requiresAuth: true,
					title: "动态 Form",
					key: "dynamicForm"
				}
			}
		]
	}
];

export default formRouter;
