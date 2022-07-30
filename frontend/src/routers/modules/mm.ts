import { RouteRecordRaw } from "vue-router";
const mmRouter: Array<RouteRecordRaw> = [
	{
		path: "/mm/pred",
		name: "mm",
		component: () => import("@/views/mm/index.vue"),
		meta: {
			requiresAuth: false,
			title: "多发性骨髓瘤细胞检测",
			key: "mm"
		}
		// redirect: "/mm/pred"
	}
];

export default mmRouter;
