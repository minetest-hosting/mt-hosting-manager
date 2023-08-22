import store from "../../store/info.js";
import CardLayout from "../layouts/CardLayout.js";

export default {
	data: () => store,
	components: {
		"card-layout": CardLayout
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-home"></i> Start
		</template>
		<div class="text-center">
			<h3>
				Minetest hosting
			</h3>
		</div>
		<a :href="'https://github.com/login/oauth/authorize?client_id=' + github_client_id + '&scope=user:email'" class="btn btn-secondary">
            <i class="fab fa-github"></i>
            Login with Github
        </a>
	</card-layout>
	`
};
