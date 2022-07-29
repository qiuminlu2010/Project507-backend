import { RouteRecordRaw } from "vue-router";
import { Layout } from "@/routers/constant";

// 首页模块
const homeRouter: Array<RouteRecordRaw> = [
	{
		path: "/",
		component: Layout,
		redirect: "/home",
		children: [
			{
				path: "/home",
				name: "home",
				component: () => import("@/views/home/index.vue"),
				meta: {
					keepAlive: true,
					requiresAuth: false,
					title: "主页",
					key: "home"
				}
			},
			{
				path: "/msg",
				name: "msg",
				component: () => import("@/views/msg/index.vue"),
				meta: {
					requiresAuth: false,
					title: "消息",
					key: "msg"
				}
			}
		]
	}
];

export default homeRouter;
