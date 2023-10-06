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
			<router-link class="btn btn-secondary" to="/pricing">
				<i class="fa fa-money-bill"></i> Pricing
			</router-link>
			&nbsp;
			<router-link class="btn btn-secondary" to="/privacy-policy">
				<i class="fa fa-section"></i> Privacy policy
			</router-link>
			&nbsp;
			<router-link class="btn btn-secondary" to="/terms-conditions">
				<i class="fa fa-section"></i> Terms and conditions
			</router-link>
			&nbsp;
			<a class="btn btn-secondary" href="mailto:hosting@minetest.ch">
				<i class="fa fa-envelope"></i> Contact
			</a>
			&nbsp;
			<a class="btn btn-secondary" href="https://github.com/minetest-hosting/mt-hosting-manager" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
			&nbsp;
			<a class="btn btn-secondary" href="https://github.com/minetest-go/mtui" target="new">
				<i class="fa-brands fa-github"></i> Powered by mtui
			</a>
			&nbsp;
			<a class="btn btn-secondary" href="https://discord.gg/Xj62fUbQkn" target="new">
				<i class="fa-brands fa-discord"></i>
				Join the discord server
			</a>
			<hr/>
			<img src="assets/minetest-hosting-600px.png" class="img img-rounded"/>
		</div>
	</card-layout>
	`
};
