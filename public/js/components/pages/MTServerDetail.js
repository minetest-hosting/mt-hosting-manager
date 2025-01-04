import CardLayout from "../layouts/CardLayout.js";
import ServerState from "../ServerState.js";
import HelpPopup from "../HelpPopup.js";
import NodeLink from "../NodeLink.js";
import ClipboardCopy from "../ClipboardCopy.js";
import ServerStatsBadge from "../ServerStatsBadge.js";
import TimestampBadge from "../TimestampBadge.js";

import { get_by_id, setup, get_latest_job, update } from "../../api/mtserver.js";
import { get_hostingdomain_suffix } from "../../service/info.js";
import { get_by_id as get_node_by_id } from "../../api/node.js";
import { get_all as get_all_backup_spaces } from "../../api/backup_space.js";
import { create as create_backup } from "../../api/backup.js";

import { has_role } from "../../service/login.js";

export default {
	props: ["id"],
	components: {
		"card-layout": CardLayout,
		"server-state": ServerState,
		"help-popup": HelpPopup,
		"node-link": NodeLink,
		"clipboard-copy": ClipboardCopy,
		"server-stats-badge": ServerStatsBadge,
		"timestamp-badge": TimestampBadge
	},
	mounted: function() {
		const server_id = this.id;
		get_all_backup_spaces().then(s => this.backup_spaces = s);
		get_by_id(server_id)
		.then(s => {
			this.server = s;
			return get_node_by_id(s.user_node_id);
		})
		.then(n => {
			this.node = n;
			this.handle = setInterval(() => this.update(), 2000);
		});
	},
	beforeUnmount: function() {
		clearInterval(this.handle);
	},
	data: function(){
		return {
			server: null,
			node: null,
			backup_spaces: [],
			backup_space: null,
			backup_scheduled: false,
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
		has_role,
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
		},
		create_backup: function() {
			this.backup_scheduled = true;
			create_backup({
				backup_space_id: this.backup_space,
				minetest_server_id: this.id
			});
		}
	},
	computed: {
		setup_running: function() {
			return (this.job && this.job.state == "RUNNING");
		},
		admin_login_url: function() {
			const s = this.server;
			return `https://${this.dns_name}/ui/api/loginadmin/${s.admin}?key=${s.jwt_key}`;
		},
		server_fresh: function() {
			// server created in the last 2 minutes
			return (Date.now() - (this.server.created*1000)) < 120000;
		},
		dns_name: function() {
			if (this.server.custom_dns_name == "") {
				return this.server.dns_name + "." + this.dns_suffix;
			} else {
				return this.server.custom_dns_name;
			}
		}
	},
	template: /*html*/`
	<card-layout title="Server details" icon="list" :breadcrumb="breadcrumb">
		<table class="table" v-if="server && node">
			<tbody>
				<tr>
					<td>ID</td>
					<td>{{server.id}}</td>
				</tr>
				<tr>
					<td>Name</td>
					<td>
						<input type="text" class="form-control" v-model="server.name" :disabled="server.state != 'RUNNING'"/>
					</td>
				</tr>
				<tr>
					<td>Created</td>
					<td>
						<timestamp-badge :timestamp="server.created" :show_duration="true"/>
					</td>
				</tr>
				<tr>
					<td>Parent node</td>
					<td>
						<node-link :node="node"/>
					</td>
				</tr>
				<tr>
					<td>Port</td>
					<td>{{server.port}}</td>
				</tr>
				<tr>
					<td>DNS Name</td>
					<td>
						<clipboard-copy :text="dns_name"/>
					</td>
				</tr>
				<tr v-if="server.state == 'RUNNING'">
					<td>
						Custom DNS Name
						<help-popup title="Custom DNS Name">
							<p>You can use your own dns name and point a CNAME record to: <b>{{node.name}}.{{dns_suffix}}</b></p>
							<p>The server will use the custom name after the UI setup has been run again</p>
						</help-popup>
					</td>
					<td>
						<input type="text" placeholder="Custom DNS Name" class="form-control" v-model="server.custom_dns_name"/>
					</td>					
				</tr>
				<tr v-if="server.state == 'RUNNING' && has_role('ADMIN')">
					<td>UI Version</td>
					<td>
						<input type="text" class="form-control" v-model="server.ui_version"/>
						<div class="alert alert-info">
							<i class="fa-solid fa-triangle-exclamation"></i>
							Don't change this, unless you <i>really</i> know what you are doing!
						</div>
					</td>
				</tr>
				<tr v-if="server.state == 'RUNNING'">
					<td>Admin-user</td>
					<td>
						<input type="text" class="form-control" v-model="server.admin"/>
					</td>
				</tr>
				<tr v-if="server.state == 'RUNNING'">
					<td>Admin login</td>
					<td>
						<i class="fa-solid fa-arrow-up-right-from-square"></i>
						<a :href="admin_login_url" target="new">{{dns_name}}/ui</a>
						<div class="alert alert-warning" v-if="server_fresh">
							<i class="fa-solid fa-triangle-exclamation"></i>
							The server was recently created, if the ui-link does not work wait another minute or two.
						</div>
					</td>
				</tr>
				<tr>
					<td>State</td>
					<td>
						<server-state :state="server.state" v-if="server.state != 'RUNNING'"/>
						<server-stats-badge v-if="server.state == 'RUNNING'" :id="server.id"/>
					</td>
				</tr>
				<tr v-if="server.state == 'RUNNING' && backup_spaces.length > 0">
					<td>Backup</td>
					<td>
						<div class="input-group">
							<select v-model="backup_space" class="form-control" :disabled="backup_scheduled">
								<option v-for="bs in backup_spaces" :value="bs.id">{{bs.name}} ({{bs.retention_days}} days retention)</option>
							</select>
							<button class="btn btn-secondary" :disabled="!backup_space || backup_scheduled" v-on:click="create_backup">
								<i class="fa fa-floppy-disk"></i>
								Create
							</button>
						</div>
						<router-link :to="'/backup_spaces/' + backup_space" v-if="backup_scheduled">
							<i class="fa fa-check"></i>
							Backup scheduled
						</router-link>
					</td>
				</tr>
				<tr v-if="server.state == 'RUNNING'">
					<td>Actions</td>
					<td>
						<div class="btn-group">
							<button class="btn btn-sm btn-success" v-on:click="save" :disabled="setup_running">
								<i class="fa fa-floppy-disk"></i>
								Save changes
							</button>
							<button class="btn btn-sm btn-outline-secondary" v-on:click="setup" :disabled="setup_running">
								<i class="fa fa-cog"></i>
								Update management UI
								<i class="fa fa-spinner fa-spin" v-if="setup_running"></i>
							</button>
							<router-link class="btn btn-sm btn-danger" :to="'/mtservers/' + server.id + '/delete'" :disabled="server.state != 'RUNNING'">
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
