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
                icon: "object-group", name: "Backup spaces", link: "backup_spaces"
            }]
		};
	},
	template: /*html*/`
	<card-layout title="Backup spaces" icon="object-group" :breadcrumb="breadcrumb">
	</card-layout>
	`
};
