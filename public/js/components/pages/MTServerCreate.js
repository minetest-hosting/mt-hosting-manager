import CardLayout from "../layouts/CardLayout.js";
import { get_hostingdomain_suffix } from "../../service/info.js";
import { create as create_server } from "../../api/mtserver.js";
import { get_all as get_all_nodes } from "../../api/node.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			user_nodes: [],
			user_node_id: "",
			port: 30000,
			name: "",
			dns_name: "",
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
			this.user_node_id = n[0].id;
			this.user_nodes = n;
		});
	},
	methods: {
		create: function() {
			create_server({
				port: this.port,
				name: this.name,
				dns_name: this.dns_name,
				user_node_id: this.user_node_id
			})
			.then(s => this.$router.push(`/mtservers/${s.id}`));
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
						<input type="number" min="1000" max="65500" class="form-control" v-model="port"/>
					</td>
				</tr>
				<tr>
					<td>DNS Prefix</td>
					<td>
						<div class="input-group">
							<input type="text" class="form-control" v-model="dns_name"/>
							<span class="input-group-text">.{{dns_suffix}}</span>
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
