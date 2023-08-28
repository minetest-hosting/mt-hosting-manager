import CardLayout from "../layouts/CardLayout.js";
import ServerState from "../ServerState.js";

import { get_by_id, setup, get_latest_job } from "../../api/mtserver.js";
import { get_hostingdomain_suffix } from "../../service/info.js";

export default {
	components: {
		"card-layout": CardLayout,
		"server-state": ServerState
	},
	mounted: function() {
		this.update();
	},
	data: function(){
		return {
			server: null,
			job: null,
			dns_suffix: get_hostingdomain_suffix(),
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "list", name: "Servers", link: "/mtservers"
			},{
				icon: "list", name: "Server detail", link: `/mtservers/${this.$route.params.id}`
			}]
		};
	},
	methods: {
		update: function() {
			get_by_id(this.$route.params.id)
			.then(s => this.server = s);

			get_latest_job(this.$route.params.id)
			.then(j => this.job = j);
		},
		setup: function() {
			setup(this.server);
		}
	},
	template: /*html*/`
	<card-layout title="Server details" icon="list" :breadcrumb="breadcrumb">
		<table class="table" v-if="server">
			<tbody>
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
				<tr>
					<td>UI Version</td>
					<td>
						<input type="text" class="form-control" v-model="server.ui_version"/>
					</td>
				</tr>
				<tr>
					<td>State</td>
					<td>
						<server-state :state="server.state"/>
					</td>
				</tr>
				<tr>
					<td>Actions</td>
					<td>
						<div class="btn-group">
							<a class="btn btn-xs btn-outline-secondary" v-on:click="setup">
								Run setup
							</a>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
	</card-layout>
	`
};
