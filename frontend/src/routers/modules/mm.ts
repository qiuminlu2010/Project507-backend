import { RouteRecordRaw } from "vue-router";
const mmRouter: Array<RouteRecordRaw> = [
	{
		path: "/autotag",
		name: "mm",
		component: () => import("@/views/mm/autotag.vue"),
		meta: {
			requiresAuth: false,
			title: "多发性骨髓瘤细胞检测",
			key: "mm"
		}
		// redirect: "/mm/pred"
	}
];

export default mmRouter;
