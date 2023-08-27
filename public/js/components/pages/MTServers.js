import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "list", name: "Servers", link: "/mtservers"
			}]
		};
	},
	template: /*html*/`
	<card-layout title="Servers" icon="list" :breadcrumb="breadcrumb">
	</card-layout>
	`
};
