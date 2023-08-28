import CardLayout from "../layouts/CardLayout.js";
import { logout, is_logged_in } from "../../service/login.js";
import { get_github_client_id } from "../../service/info.js";

export default {
	data: function() {
		return {
			github_client_id: get_github_client_id(),
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "user", name: "Login", link: "/login"
            }]
		};
	},
	components: {
		"card-layout": CardLayout
	},
	methods: {
		logout: logout
	},
	computed: {
		is_logged_in: is_logged_in
	},
	template: /*html*/`
	<card-layout title="Login" icon="user" :breadcrumb="breadcrumb">
		<a :href="'https://github.com/login/oauth/authorize?client_id=' + github_client_id + '&scope=user:email'" class="btn btn-secondary" v-if="!is_logged_in">
			<i class="fab fa-github"></i>
			Login with Github
		</a>
		<a class="btn btn-primary" v-on:click="logout" v-if="is_logged_in">
			Logout
		</a>
	</card-layout>
    `
};
