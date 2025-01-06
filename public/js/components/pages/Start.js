import CardLayout from "../layouts/CardLayout.js";
import Breadcrumb from "../Breadcrumb.js";

export default {
	components: {
		"card-layout": CardLayout,
		"bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [{icon: "home", name: "Home", link: "/"}]
		};
	},
	template: /*html*/`
	<card-layout title="Home" icon="home" :breadcrumb="breadcrumb">
		<div class="text-center">
			<h4>Minetest hosting</h4>
			<img src="assets/minetest-hosting-80px.png" class="img img-rounded"/>
			<hr/>
			<router-link class="btn btn-secondary" to="/pricing">
				<i class="fa fa-money-bill"></i> Pricing
			</router-link>
			&nbsp;
			<a class="btn btn-secondary" href="mailto:hosting@minetest.ch">
				<i class="fa fa-envelope"></i> Contact
			</a>
			&nbsp;
			<a class="btn btn-secondary" href="https://discord.gg/Xj62fUbQkn" target="new">
				<i class="fa-brands fa-discord"></i>
				Join the discord server
			</a>
			&nbsp;
			<a class="btn btn-outline-secondary" href="https://github.com/minetest-hosting/mt-hosting-manager" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
			&nbsp;
			<a class="btn btn-outline-secondary" href="https://github.com/minetest-go/mtui" target="new">
				<i class="fa-brands fa-github"></i> Powered by mtui
			</a>
		</div>
	</card-layout>
	`
};
