<template>
	<el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" @keyup.enter="onShowVcode" size="large">
		<div class="login-logo">
			<span class="logo-text">你好，朋友!</span>
		</div>
		<el-form-item prop="username">
			<el-input v-model="loginForm.username" placeholder="用户名">
				<template #prefix>
					<el-icon class="el-input__icon"><user /></el-icon>
				</template>
			</el-input>
		</el-form-item>
		<el-form-item prop="password">
			<el-input type="password" show-password v-model="loginForm.password" placeholder="密码">
				<template #prefix>
					<el-icon class="el-input__icon"><lock /></el-icon>
				</template>
			</el-input>
		</el-form-item>
		<el-form-item prop="password2">
			<el-input type="password2" show-password v-model="loginForm.password2" placeholder="再次输入密码">
				<template #prefix>
					<el-icon class="el-input__icon"><lock /></el-icon>
				</template>
			</el-input>
		</el-form-item>
		<!-- <el-form-item prop="email">
			<el-input type="text" clearable v-model="loginForm.email" autofocus placeholder="请输入登录邮箱">
				<template #prefix>
					<el-icon i="ep-user"></el-icon>
				</template>
			</el-input>
		</el-form-item> -->

		<el-form-item prop="code">
			<el-input v-model="loginForm.code" placeholder="请输入验证码">
				<template #prefix>
					<el-icon class="el-input__icon"><clock /></el-icon>
				</template>
				<template #append>
					<el-button type="primary" @click="sendCode(loginFormRef)">
						{{ codeText }}
					</el-button>
				</template>
			</el-input>
		</el-form-item>

		<el-form-item class="login-btn">
			<!-- <el-button :icon="CircleClose" round @click="resetForm(loginFormRef)" size="large">注册</el-button> -->
			<el-button :icon="UserFilled" round @click="onShowVcode" size="large" type="primary" :loading="loading"> 注册 </el-button>
		</el-form-item>
	</el-form>
	<Vcode :show="isShowVcode" @success="signup(loginFormRef)" @close="onCloseVcode" />
</template>

<script setup lang="ts">
import { ref, reactive, inject, computed } from "vue";
import { Login } from "@/api/interface";
import { InjectProps } from "../interface/index";
import { UserFilled } from "@element-plus/icons-vue";
import { ElForm, ElNotification } from "element-plus";
import { ElMessage } from "element-plus";
import { signupApi } from "@/api/modules/login";
import { validatorPassword, validatorUsername } from "../utils";
import md5 from "js-md5";
import Vcode from "vue3-puzzle-vcode";

// const globalStore = GlobalStore();

// 定义 formRef（校验规则）
type FormInstance = InstanceType<typeof ElForm>;
const loginFormRef = ref<FormInstance>();
// 验证再次输入密码
const validatePass2 = (rule: any, value: any, callback: any) => {
	if (value === "") {
		callback(new Error("请再次输入密码"));
	} else if (value !== loginForm.password) {
		callback(new Error("两次密码不一致"));
	} else {
		callback();
	}
};

const loginRules = reactive({
	username: [
		{ required: true, message: "请输入用户名", trigger: "blur" },
		{ validator: validatorUsername, trigger: ["blur", "change"] }
	],
	password: [
		{ required: true, message: "请输入密码", trigger: "blur" },
		{ validator: validatorPassword, trigger: ["blur", "change"] }
	],
	password2: [
		{
			required: true,
			validator: validatePass2,
			trigger: ["blur", "change"]
		}
	]
});

// 登录表单数据
const loginForm = reactive({
	username: "",
	password: "",
	password2: "",
	email: "",
	code: ""
});
const loading = ref<boolean>(false);
const isShowVcode = ref<boolean>(false);
const onShowVcode = () => {
	isShowVcode.value = true;
};
const onCloseVcode = () => {
	isShowVcode.value = false;
};

// const router = useRouter();
// login
const signup = (formEl: FormInstance | undefined) => {
	if (!formEl) return;
	formEl.validate(valid => {
		if (valid) {
			loading.value = true;
			try {
				let requestLoginForm: Login.ReqLoginForm = {
					username: loginForm.username,
					password: md5(loginForm.password)
				};
				// eslint-disable-next-line @typescript-eslint/no-unused-vars
				const res = signupApi(requestLoginForm);
				// const res = await signupApi(requestLoginForm);
				// console.log("响应数据", res.data);
				// globalStore.setToken(res.data!.access_token);
				ElMessage.success("注册成功！");
				isLogin.value = !isLogin.value;
			} finally {
				loading.value = false;
				isShowVcode.value = false;
			}
		} else {
			ElMessage.error("注册失败！");
			loading.value = false;
			isShowVcode.value = false;
			return false;
		}
	});
};
// 登录时间戳
// const timestamp = ref("");
// 验证码数据
const codeDatas = reactive({
	timer: null,
	num: 60
});
// 验证码倒计时
const codeText = computed(() => {
	return codeDatas.num !== 60 ? `剩余${codeDatas.num}秒` : "获取验证码";
});

// 获取验证码
function sendCode(formEl: InstanceType<typeof ElForm> | undefined) {
	if (!formEl || codeDatas.num !== 60) return;
	formEl.validateField("email", bool => {
		if (bool) {
			codeDatas.num--;
			codeDatas.timer = setInterval(() => {
				if (codeDatas.num < 1) {
					if (codeDatas.timer != null) {
						clearInterval(codeDatas.timer);
						codeDatas.timer = null;
					}
					codeDatas.num = 60;
				} else {
					codeDatas.num--;
				}
			}, 1000);
			ElNotification.success("发送验证码成功");
			loginForm.code = "123456";
			// timestamp.value = res?.data?.timestamp;
			// else ElNotification.success("发送验证码失败");
		}
	});
}
// * 以下数据都为自己测试使用，不参与功能开发
// inject
const isLogin = inject("isLogin") as InjectProps;
</script>

<style scoped lang="scss">
@import "../index.scss";
</style>
