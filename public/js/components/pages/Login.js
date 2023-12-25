import CardLayout from "../layouts/CardLayout.js";
import { logout, is_logged_in, login } from "../../service/login.js";
import { get_github_login, get_discord_login, get_mesehub_login, get_cdb_login } from "../../service/info.js";

export default {
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "user", name: "Login", link: "/login"
            }],
			username: "",
			password: "",
			login_error: false
		};
	},
	components: {
		"card-layout": CardLayout
	},
	methods: {
		logout,
		login: function() {
			login({
				username: this.username,
				password: this.password
			})
			.then(success => {
				if (!success) {
					this.login_error = true;
				} else {
					this.$router.push("/profile");
				}
			});
		}
	},
	computed: {
		is_logged_in,
		get_github_login,
		get_discord_login,
		get_mesehub_login,
		get_cdb_login
	},
	template: /*html*/`
	<card-layout title="Login" icon="user" :breadcrumb="breadcrumb">
		<div class="row" v-if="!is_logged_in">
			<div class="col-6">
				<h5>Login with username and password</h5>
				<form @submit.prevent="login">
					<input type="text" class="form-control" placeholder="Username" v-model="username" v-bind:class="{'is-invalid': login_error}"/>
					&nbsp;
					<input type="password" class="form-control" placeholder="Password" v-model="password" v-bind:class="{'is-invalid': login_error}"/>
					&nbsp;
					<button type="submit" class="btn btn-primary w-100">
						<i class="fa-solid fa-right-to-bracket"></i>
						Login
					</button>
				</form>
			</div>
			<div class="col-6">
				<h5>Login with external provider</h5>
				<a :href="get_github_login" class="btn btn-secondary w-100" v-if="get_github_login">
					<i class="fab fa-github"></i>
					Login with Github
				</a>
				&nbsp;
				<a :href="get_discord_login" class="btn btn-secondary w-100" v-if="get_discord_login">
					<i class="fab fa-discord"></i>
					Login with Discord
				</a>
				&nbsp;
				<a :href="get_mesehub_login" class="btn btn-secondary w-100" v-if="get_mesehub_login">
					<img src="assets/default_mese_crystal.png">
					Login with Mesehub
				</a>
				&nbsp;
				<a :href="get_cdb_login" class="btn btn-secondary w-100" v-if="get_cdb_login">
					<img src="assets/contentdb.png" height="20" width="20">
					Login with ContentDB
				</a>
			</div>
		</div>
		<hr>
		<div class="row" v-if="!is_logged_in">
			<div class="col-6">
				Register a new account <router-link to="/register">here</router-link>
			</div>
			<div class="col6">
			</div>
		</div>
		<a class="btn btn-primary w-100" v-on:click="logout" v-if="is_logged_in">
			<i class="fa-solid fa-right-from-bracket"></i>
			Logout
		</a>
	</card-layout>
    `
};
