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
			<hr/>
			<img src="assets/minetest-hosting-600px.png" class="img img-rounded"/>
		</div>
	</card-layout>
	`
};
