import { defineStore } from "pinia";
import { GlobalState, ThemeConfigProp } from "./interface";
import { createPinia } from "pinia";
import piniaPersist from "pinia-plugin-persist";
import piniaPersistConfig from "@/config/piniaPersist";

// defineStore 调用后返回一个函数，调用该函数获得 Store 实体
export const GlobalStore = defineStore({
	// id: 必须的，在所有 Store 中唯一
	id: "GlobalState",
	// state: 返回对象的函数
	state: (): GlobalState => ({
		// token
		token: "",
		// uuid
		uuid: "",
		expireTime: 0,
		uid: 0,
		username: "",
		name: "",
		avatar: "",
		// element组件大小
		assemblySize: "default",
		// language
		language: "",
		// themeConfig
		themeConfig: {
			// 默认 primary 主题颜色
			primary: "#409eff",
			// 是否开启深色模式
			isDark: false
		}
	}),
	getters: {},
	actions: {
		// setToken
		setToken(token: string) {
			this.token = token;
		},
		setUUID(uuid: string) {
			this.uuid = uuid;
		},
		setExpireTime(time: number) {
			this.expireTime = time;
		},
		setUID(uid: number) {
			this.uid = uid;
		},
		setUsername(username: string) {
			this.username = username;
		},
		setName(name: string | undefined) {
			this.name = name;
		},
		setAvatar(avatar: string | undefined) {
			this.avatar = avatar;
		},
		setAssemblySizeSize(assemblySize: string) {
			this.assemblySize = assemblySize;
		},
		// updateLanguage
		updateLanguage(language: string) {
			this.language = language;
		},
		// setThemeConfig
		setThemeConfig(themeConfig: ThemeConfigProp) {
			this.themeConfig = themeConfig;
		}
	},
	persist: piniaPersistConfig("GlobalState")
});

// export const CommentStore = defineStore({
// 	id: "CommentState",
// 	state: (): CommentState => ({
// 		currentCommentList: [],
// 		editor: null,
// 		scrollBar: null,
// 		emojiList: []
// 	}),
// 	getters: {},
// 	actions: {
// 		setCurrentCommentList(CommentList: any) {
// 			this.currentCommentList = CommentList;
// 		},
// 		handleEditorCreated(editor: any) {
// 			this.editor = editor;
// 		},
// 		handleEditorChange() {
// 			console.log("change:", this!.editor.getHtml().trim());
// 		}
// 	}
// });

// piniaPersist(持久化)
const pinia = createPinia();
pinia.use(piniaPersist);

export default pinia;
