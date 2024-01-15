import CardLayout from "../layouts/CardLayout.js";

import { get_hostingdomain_suffix } from "../../service/info.js";
import { create as create_server, create_validate } from "../../api/mtserver.js";
import { get_all as get_all_nodes } from "../../api/node.js";
import random_name from "../../util/random_name.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		const n = random_name().toLowerCase().replaceAll("-", "");
		return {
			validation_result: {},
			user_nodes: [],
			user_node_id: this.$route.query.node ? this.$route.query.node : "",
			port: 30000,
			admin: "admin",
			name: n,
			dns_name: n,
			dns_suffix: get_hostingdomain_suffix(),
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "list", name: "Servers", link: "/mtservers"
			},{
				icon: "plus", name: "Create server", link: "/mtservers/create"
			}]
		};
	},
	mounted: function() {
		get_all_nodes()
		.then(n => {
			const nodelist = n.filter(node => node.state == "RUNNING");
			if (this.user_node_id == "" && nodelist.length) {
				// default node id
				this.user_node_id = nodelist[0].id;
			}
			this.user_nodes = nodelist;
		});
	},
	methods: {
		create: function() {
			const server = {
				port: this.port,
				name: this.name,
				dns_name: this.dns_name,
				admin: this.admin,
				user_node_id: this.user_node_id
			};

			create_validate(server)
			.then(v => {
				if (v.valid) {
					return create_server(server)
					.then(s => this.$router.push(`/mtservers/${s.id}`));
				}

				this.validation_result = v;
			});
		}
	},
	computed: {
		valid: function() {
			return this.port && this.name && this.dns_name && this.user_node_id;
		}
	},
	template: /*html*/`
	<card-layout title="Create server" icon="plus" :breadcrumb="breadcrumb">
		<table class="table">
			<tbody>
				<tr>
					<td>Node</td>
					<td>
						<select v-model="user_node_id" class="form-control">
							<option v-for="node in user_nodes" :value="node.id">{{node.name}}</option>
						</select>
					</td>
				</tr>
				<tr>
					<td>Name</td>
					<td>
						<input type="text" class="form-control" v-model="name"/>
					</td>
				</tr>
				<tr>
					<td>Port</td>
					<td>
						<input type="number" min="1000" max="65500" class="form-control" v-bind:class="{'is-invalid':validation_result.port_invalid || validation_result.port_used}" v-model="port"/>
						<div class="invalid-feedback" v-if="validation_result.port_invalid">
							Port number is invalid
						</div>
						<div class="invalid-feedback" v-if="validation_result.port_used">
							Port number already used
						</div>
					</td>
				</tr>
				<tr>
					<td>Admin-user</td>
					<td>
						<input type="text" class="form-control" v-bind:class="{'is-invalid':validation_result.admin_name_invalid}" v-model="admin"/>
						<div class="invalid-feedback" v-if="validation_result.admin_name_invalid">
							Username invalid
						</div>
					</td>
				</tr>
				<tr>
					<td>DNS Prefix</td>
					<td>
						<div class="input-group">
							<input type="text" class="form-control" v-bind:class="{'is-invalid':validation_result.server_name_invalid||validation_result.server_name_used}" v-model="dns_name"/>
							<span class="input-group-text">.{{dns_suffix}}</span>
							<div class="invalid-feedback" v-if="validation_result.server_name_invalid">
								Servername invalid, only the letters a to z and numbers 0 to 9 can be used
							</div>
							<div class="invalid-feedback" v-if="validation_result.server_name_used">
								Servername already used
							</div>
							<div class="invalid-feedback" v-if="validation_result.server_name_too_short">
								Servername too short (min 5 letters)
							</div>
							<div class="invalid-feedback" v-if="validation_result.server_name_reserved">
								Servername is reserved
							</div>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
		<div class="row">
			<div class="col-12">
				<button class="btn btn-success w-100" :disabled="!valid" v-on:click="create">
					<i class="fa fa-plus"></i>
					Create server
				</button>
			</div>
		</div>
	</card-layout>
	`
};
