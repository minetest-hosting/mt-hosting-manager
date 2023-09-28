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
			<a class="btn btn-secondary" href="https://github.com/minetest-hosting/mt-hosting-manager" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
			&nbsp;
			<a class="btn btn-secondary" href="https://github.com/minetest-go/mtui" target="new">
				<i class="fa-brands fa-github"></i> Powered by mtui
			</a>
			<hr/>
			<img src="assets/minetest-hosting-600px.png" class="img img-rounded"/>
		</div>
	</card-layout>
	`
};
