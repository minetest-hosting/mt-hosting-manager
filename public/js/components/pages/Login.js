import CardLayout from "../layouts/CardLayout.js";
import { logout } from "../../api/login.js";
import { get_github_client_id } from "../../service/info.js";

export default {
	data: function() {
		return {
			github_client_id: get_github_client_id()
		};
	},
	components: {
		"card-layout": CardLayout
	},
	methods: {
		logout: logout
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-user"></i> Login
		</template>
		<a :href="'https://github.com/login/oauth/authorize?client_id=' + github_client_id + '&scope=user:email'" class="btn btn-secondary">
			<i class="fab fa-github"></i>
			Login with Github
		</a>
		<a class="btn btn-primary" v-on:click="logout">
			Logout
		</a>
	</card-layout>
    `
};
