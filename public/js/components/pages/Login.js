import CardLayout from "../layouts/CardLayout.js";
import { logout, is_logged_in } from "../../service/login.js";
import { get_github_client_id, get_discord_client_id, get_mesehub_client_id, get_baseurl } from "../../service/info.js";

export default {
	data: function() {
		return {
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
		logout,
		github_href: function() {
			return `https://github.com/login/oauth/authorize?client_id=${get_github_client_id()}&scope=user:email`;
		},
		discord_href: function() {
			return `https://discord.com/api/oauth2/authorize?client_id=${get_discord_client_id()}&redirect_uri=${encodeURIComponent(get_baseurl()+'/oauth_callback/discord')}&response_type=code&scope=identify%20email`;
		},
		mesehub_href: function() {
			return `https://git.minetest.land/login/oauth/authorize?client_id=${get_mesehub_client_id()}&redirect_uri=${encodeURIComponent(get_baseurl()+'/oauth_callback/mesehub')}&response_type=code&state=STATE&scope=email`;
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
		<a :href="github_href()" class="btn btn-secondary" v-if="!is_logged_in && get_github_client_id">
			<i class="fab fa-github"></i>
			Login with Github
		</a>
		&nbsp;
		<a :href="discord_href()" class="btn btn-secondary" v-if="!is_logged_in && get_discord_client_id">
			<i class="fab fa-discord"></i>
			Login with Discord
		</a>
		&nbsp;
		<a :href="mesehub_href()" class="btn btn-secondary" v-if="!is_logged_in && get_mesehub_client_id">
			<img src="assets/default_mese_crystal.png">
			Login with Mesehub
		</a>
		&nbsp;
		<a class="btn btn-primary" v-on:click="logout" v-if="is_logged_in">
			Logout
		</a>
	</card-layout>
    `
};
