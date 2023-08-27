import CardLayout from "../layouts/CardLayout.js";
import { get_by_id } from "../../api/mtserver.js";
import { get_hostingdomain_suffix } from "../../service/info.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	mounted: function() {
		get_by_id(this.$route.params.id)
		.then(s => this.server = s);
	},
	data: function(){
		return {
			server: null,
			dns_suffix: get_hostingdomain_suffix()
		};
	},
	template: /*html*/`
	<card-layout title="Server details" icon="list">
		<table class="table" v-if="server">
			<tr>
				<td>ID</td>
				<td>{{server.id}}</td>
			</tr>
			<tr>
				<td>Name</td>
				<td>{{server.name}}</td>
			</tr>
			<tr>
				<td>Port</td>
				<td>{{server.port}}</td>
			</tr>
			<tr>
				<td>DNS Name</td>
				<td>
					<i class="fa-solid fa-arrow-up-right-from-square"></i>
					<a :href="'https://' + server.dns_name + '.' + dns_suffix">{{server.dns_name}}.{{dns_suffix}}</a>
				</td>
			</tr>
		</table>
	</card-layout>
	`
};
