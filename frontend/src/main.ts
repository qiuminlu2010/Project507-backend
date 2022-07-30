import { createApp } from "vue";
import App from "./App.vue";
// reset style sheet
import "@/styles/reset.scss";
// CSS common style sheet
import "@/styles/common.scss";
// iconfont css
import "@/assets/iconfont/iconfont.scss";
// font css
import "@/assets/fonts/font.scss";
// element plus
import ElementPlus from "element-plus";
// element icons
import * as Icons from "@element-plus/icons-vue";
// element css
import "element-plus/dist/index.css";
// element dark
import "element-plus/theme-chalk/dark/css-vars.css";
// custom element dark
import "@/styles/element-dark.scss";
// custom element css
import "@/styles/element.scss";

// custom directives
import directives from "@/directives/index";
// vue Router
import router from "@/routers/index";
// i18n
import I18n from "@/language/index";
// pinia store
import pinia from "@/store/index";
import "uno.css";

//编辑器
import "@wangeditor/editor/dist/css/style.css";
//安全html
import VueDOMPurifyHTML from "vue-dompurify-html";
//日期时间库
import moment from "moment";
import "moment/dist/locale/zh-cn";
//瀑布流
import V3waterfall from "v3-waterfall";
import "v3-waterfall/style.css";
moment.locale("zh-cn");
//视频播放器
import vue3videoPlay from "vue3-video-play"; // 引入组件
import "vue3-video-play/dist/style.css"; // 引入css
const app = createApp(App);
// 注册element Icons组件
Object.keys(Icons).forEach(key => {
	app.component(key, Icons[key as keyof typeof Icons]);
});

app
	.use(router)
	.use(I18n)
	.use(pinia)
	.use(directives)
	.use(ElementPlus)
	.use(VueDOMPurifyHTML)
	.use(V3waterfall)
	.use(vue3videoPlay)
	.mount("#app");
