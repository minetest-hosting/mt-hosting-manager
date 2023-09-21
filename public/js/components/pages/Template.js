import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			breadcrumb: [{icon: "home", name: "Home", link: "/"}]
		};
	},
	template: /*html*/`
	<card-layout title="Template" icon="home" :breadcrumb="breadcrumb">
	</card-layout>
	`
};
