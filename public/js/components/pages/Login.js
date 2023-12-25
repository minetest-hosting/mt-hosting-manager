import CardLayout from "../layouts/CardLayout.js";
import { logout, is_logged_in, login } from "../../service/login.js";
import { get_github_client_id, get_discord_client_id, get_mesehub_client_id, get_baseurl } from "../../service/info.js";

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
		github_href: function() {
			return `https://github.com/login/oauth/authorize?client_id=${get_github_client_id()}&scope=user:email`;
		},
		discord_href: function() {
			return `https://discord.com/api/oauth2/authorize?client_id=${get_discord_client_id()}&redirect_uri=${encodeURIComponent(get_baseurl()+'/oauth_callback/discord')}&response_type=code&scope=identify%20email`;
		},
		mesehub_href: function() {
			return `https://git.minetest.land/login/oauth/authorize?client_id=${get_mesehub_client_id()}&redirect_uri=${encodeURIComponent(get_baseurl()+'/oauth_callback/mesehub')}&response_type=code&state=STATE&scope=email`;
		},
		login: function() {
			login({
				username: this.username,
				password: this.password
			})
			.then(success => {
				if (!success) {
					this.login_error = true;
				}
			});
		}
	},
	computed: {
		is_logged_in,
		get_github_client_id,
		get_discord_client_id,
		get_mesehub_client_id
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
				<a :href="github_href()" class="btn btn-secondary w-100" v-if="get_github_client_id">
					<i class="fab fa-github"></i>
					Login with Github
				</a>
				&nbsp;
				<a :href="discord_href()" class="btn btn-secondary w-100" v-if="get_discord_client_id">
					<i class="fab fa-discord"></i>
					Login with Discord
				</a>
				&nbsp;
				<a :href="mesehub_href()" class="btn btn-secondary w-100" v-if="get_mesehub_client_id">
					<img src="assets/default_mese_crystal.png">
					Login with Mesehub
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
