import CardLayout from "../layouts/CardLayout.js";
import ServerList from "../ServerList.js";

import { get_all as get_all_servers } from "../../api/mtserver.js";

export default {
	components: {
		"card-layout": CardLayout,
		"server-list": ServerList
	},
	data: function() {
		return {
			servers: [],
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "list", name: "Servers", link: "/mtservers"
			}]
		};
	},
	mounted: function() {
		get_all_servers().then(s => this.servers = s);
	},
	template: /*html*/`
	<card-layout title="Servers" icon="list" :breadcrumb="breadcrumb">
		<server-list :list="servers"/>
	</card-layout>
	`
};
