<template>
	<div class="login-register">
		<div class="contain">
			<div class="big-box" :class="{ active: isLogin }">
				<div class="big-contain" key="bigContainLogin" v-if="isLogin">
					<div class="btitle">账户登录</div>
					<div class="bform">
						<input type="email" placeholder="邮箱" v-model="form.useremail" />
						<span class="errTips" v-if="emailError">* 邮箱填写错误 *</span>
						<input type="password" placeholder="密码" v-model="form.userpwd" />
						<span class="errTips" v-if="emailError">* 密码填写错误 *</span>
					</div>
					<button class="bbutton" @click="login">登录</button>
				</div>
				<div class="big-contain" key="bigContainRegister" v-else>
					<div class="btitle">创建账户</div>
					<div class="bform">
						<input type="text" placeholder="用户名" v-model="form.username" />
						<span class="errTips" v-if="existed">* 用户名已经存在！ *</span>
						<input type="email" placeholder="邮箱" v-model="form.useremail" />
						<input type="password" placeholder="密码" v-model="form.userpwd" />
					</div>
					<button class="bbutton" @click="register">注册</button>
				</div>
			</div>
			<div class="small-box" :class="{ active: isLogin }">
				<div class="small-contain" key="smallContainRegister" v-if="isLogin">
					<div class="stitle">你好，朋友!</div>
					<p class="scontent">开始注册，和我们一起旅行</p>
					<button class="sbutton" @click="changeType">注册</button>
				</div>
				<div class="small-contain" key="smallContainLogin" v-else>
					<div class="stitle">欢迎回来!</div>
					<p class="scontent">与我们保持联系，请登录你的账户</p>
					<button class="sbutton" @click="changeType">登录</button>
				</div>
			</div>
		</div>
	</div>
</template>
<script>
export default {
	name: "LoginRegister",
	data() {
		return {
			isLogin: false,
			emailError: false,
			passwordError: false,
			existed: false,
			form: {
				username: "",
				useremail: "",
				userpwd: ""
			}
		};
	},
	methods: {
		changeType() {
			this.isLogin = !this.isLogin;
			this.form.username = "";
			this.form.useremail = "";
			this.form.userpwd = "";
		}
	}
};
</script>
