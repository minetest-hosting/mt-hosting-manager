import CardLayout from "../layouts/CardLayout.js";
import ServerState from "../ServerState.js";

import { get_by_id, setup, get_latest_job, update } from "../../api/mtserver.js";
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
			.then(j => {
				j.state = "RUNNING";
				this.job = j;
			});
		},
		save: function() {
			update(this.server);
		}
	},
	computed: {
		setup_running: function() {
			return (this.job && this.job.state == "RUNNING");
		},
		admin_login_url: function() {
			const s = this.server;
			return `https://${s.dns_name}.${this.dns_suffix}/ui/api/loginadmin/${s.admin}?key=${s.jwt_key}`;
		},
		server_fresh: function() {
			// server created in the last 2 minutes
			return (Date.now() - (this.server.created*1000)) < 120000;
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
					<td>
						<input type="text" class="form-control" v-model="server.name"/>
					</td>
				</tr>
				<tr>
					<td>Port</td>
					<td>{{server.port}}</td>
				</tr>
				<tr>
					<td>DNS Name</td>
					<td>
						{{server.dns_name}}.{{dns_suffix}}
					</td>
				</tr>
				<tr>
					<td>Admin login</td>
					<td>
						<i class="fa-solid fa-arrow-up-right-from-square"></i>
						<a :href="admin_login_url" target="new">{{server.dns_name}}.{{dns_suffix}}/ui</a>
						<div class="alert alert-warning" v-if="server_fresh">
							<i class="fa-solid fa-triangle-exclamation"></i>
							The server was recently created, if the ui-link does not work yet wait another minute or two.
						</div>
					</td>
				</tr>
				<tr>
					<td>UI Version</td>
					<td>
						<input type="text" class="form-control" v-model="server.ui_version"/>
						<div class="alert alert-info">
							<i class="fa-solid fa-triangle-exclamation"></i>
							Don't change this, unless you <i>really</i> know what you are doing!
						</div>
					</td>
				</tr>
				<tr>
					<td>Admin-user</td>
					<td>
						<input type="text" class="form-control" v-model="server.admin"/>
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
							<button class="btn btn-xs btn-success" v-on:click="save" :disabled="setup_running">
								<i class="fa fa-floppy-disk"></i>
								Save changes
							</button>
							<button class="btn btn-xs btn-outline-secondary" v-on:click="setup" :disabled="setup_running">
								<i class="fa fa-cog"></i>
								Update management UI
								<i class="fa fa-spinner fa-spin" v-if="setup_running"></i>
							</button>
							<router-link class="btn btn-xs btn-danger" :to="'/mtservers/' + server.id + '/delete'" :disabled="setup_running">
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
