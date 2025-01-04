import CardLayout from "../layouts/CardLayout.js";
import ServerList from "../ServerList.js";

import { search } from "../../api/mtserver.js";

export default {
	components: {
		"card-layout": CardLayout,
		"server-list": ServerList
	},
	data: function() {
		return {
			show_archived: false,
			servers: [],
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "list", name: "Servers", link: "/mtservers"
			}]
		};
	},
	watch: {
		show_archived: function() {
			this.update();
		}
	},
	mounted: function() {
		this.update();
	},
	methods: {
		update: async function() {
			let s = {};
			if (!this.show_archived) {
				// limit search to active nodes
				s.state = "RUNNING";
			}
			this.servers = await search(s);
		}
	},
	template: /*html*/`
	<card-layout title="Servers" icon="list" :breadcrumb="breadcrumb">
		<div class="form-check">
			<input class="form-check-input" type="checkbox" v-model="show_archived" value="" id="show_archived">
			<label class="form-check-label" for="show_archived">
				Show archived servers
			</label>
		</div>
		<server-list :list="servers" :show_parent="true"/>
	</card-layout>
	`
};
