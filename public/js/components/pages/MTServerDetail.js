import CardLayout from "../layouts/CardLayout.js";
import ServerState from "../ServerState.js";

import defer from "../../util/defer.js";

import { get_by_id, setup, get_latest_job } from "../../api/mtserver.js";
import { get_hostingdomain_suffix } from "../../service/info.js";

export default {
	props: ["id"],
	components: {
		"card-layout": CardLayout,
		"server-state": ServerState
	},
	mounted: function() {
		const server_id = this.id;
		get_by_id(server_id)
		.then(s => this.server = s);

		this.handle = setInterval(() => this.update(), 2000);
	},
	beforeUnmount: function() {
		clearInterval(this.handle);
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
				icon: "list", name: "Server detail", link: `/mtservers/${this.id}`
			}]
		};
	},
	methods: {
		update: function() {
			const server_id = this.id;

			get_latest_job(server_id)
			.then(j => this.job = j);

			// update state
			get_by_id(server_id)
			.then(s => this.server.state = s.state);
		},
		setup: function() {
			setup(this.server)
			.then(() => defer(500))
			.then(() => this.update());
		}
	},
	computed: {
		setup_running: function() {
			return (this.job && this.job.state == "RUNNING");
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
						<a :href="'https://' + server.dns_name + '.' + dns_suffix" target="new">{{server.dns_name}}.{{dns_suffix}}</a>
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
							<button class="btn btn-xs btn-outline-secondary" v-on:click="setup" :disabled="setup_running">
								<i class="fa fa-cog"></i>
								Run setup
								<i class="fa fa-spinner fa-spin" v-if="setup_running"></i>
							</button>
							<router-link class="btn btn-xs btn-danger" :to="'/mtservers/' + server.id + '/delete'">
								<i class="fa fa-trash"></i>
								Delete
							</router-link>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
	</card-layout>
	`
};
