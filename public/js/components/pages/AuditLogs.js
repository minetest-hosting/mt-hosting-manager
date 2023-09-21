import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			}, {
				icon: "rectangle-list", name: "Audit-Logs", link: "/audit-logs"
			}]
		};
	},
	template: /*html*/`
	<card-layout title="Audit-Logs" icon="rectangle-list" :breadcrumb="breadcrumb" fullwidth="true">
		<div class="row">
			<div class="col-2">
			</div>
			<div class="col-2">
			</div>
			<div class="col-2">
			</div>
			<div class="col-2">
			</div>
		</div>
	</card-layout>
	`
};
