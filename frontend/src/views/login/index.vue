<template>
	<div class="login-container flx-center">
		<SwitchDark class="dark"></SwitchDark>
		<el-row class="login-content" type="flex" align="middle" justify="center">
			<el-col :xs="20" :sm="20" :md="16" :lg="8" :xl="8">
				<div class="box">
					<div class="small-box" :class="{ active: isLogin }">
						<div class="small-contain" key="smallContainRegister" v-if="isLogin">
							<img class="login-icon" src="@/assets/images/logo.svg" alt="" />
							<el-button class="sbutton" round @click="changeType" size="large">登录</el-button>
							<el-button class="sbutton" round size="large">游客</el-button>
						</div>
						<div class="small-contain" key="smallContainLogin" v-else>
							<img class="login-icon" src="@/assets/images/logo.svg" alt="" />
							<el-button class="sbutton" round @click="changeType" size="large">注册</el-button>
							<el-button class="sbutton" round size="large">游客</el-button>
						</div>
					</div>
					<div class="big-box" :class="{ active: isLogin }">
						<div class="login-form register-box" v-if="isLogin">
							<RegisterForm ref="loginRef"></RegisterForm>
						</div>
						<div class="login-form login-box" v-else>
							<LoginForm ref="loginRef" :age="'20'" :address="['天府三街', '天府四街']" :obj="obj"></LoginForm>
						</div>
					</div>
				</div>
			</el-col>
		</el-row>

		<!-- <div class="square">
			<ul>
				<li></li>
				<li></li>
				<li></li>
				<li></li>
				<li></li>
				<li></li>
				<li></li>
			</ul>
		</div>
		<div class="circle">
			<ul>
				<li></li>
				<li></li>
				<li></li>
				<li></li>
				<li></li>
				<li></li>
				<li></li>
			</ul>
		</div> -->
	</div>
</template>

<script setup lang="ts" name="login">
import SwitchDark from "@/components/SwitchDark/index.vue";
import LoginForm from "./components/LoginForm.vue";
import RegisterForm from "./components/RegisterForm.vue";
import { LoginFormExpose } from "./interface/index";
import { Login } from "@/api/interface/index";
import { ref, reactive, provide } from "vue";

// 传给子组件方法（submitParent）
// const submitParent = (LoginFrom: Login.ReqLoginForm) => {
// 	console.log(LoginFrom);
// };

// 使用子组件数据
const loginRef = ref<LoginFormExpose | null>(null);
// const consoleRef = () => {
// 	console.log(loginRef.value?.count);
// 	loginRef.value?.consoleNumber("liutao");
// };

// 定义数据类型两种方法
// const item: Ref<string> = ref("111");
// const item1 = ref<number>(111);

const obj = reactive<Login.ReqLoginForm>({
	username: "admin",
	password: "123456"
});
// 使用toRefs解构数据（在template中直接用）
// const { username, password } = toRefs(obj);
const isLogin = ref<boolean>(false);
const changeType = () => {
	isLogin.value = !isLogin.value;
};

// * 以下数据都为自己测试使用，不参与功能开发
// 使用provide传方法
let phone = ref<string>("iphone");
// provide
provide("provideState", {
	name: "liutao",
	age: "20",
	changeName: () => {
		// console.log(phone.value);
		setTimeout(() => {
			phone.value = "huawei";
			// console.log(phone.value);
		}, 1000);
	}
});
provide("isLogin", isLogin);
</script>

<style scoped lang="scss">
@import "./index.scss";
</style>
