import CardLayout from "../layouts/CardLayout.js";
import ServerLink from "../ServerLink.js";

import { get_by_id, update as update_node } from "../../api/node.js";
import { get_hostingdomain_suffix } from "../../service/info.js";
import { get_all as get_all_servers } from "../../api/mtserver.js";

import format_time from "../../util/format_time.js";

const bytes_in_gb = 1000 * 1000 * 1000;

function get_gb_rounded(bytes) {
	return parseInt(bytes / bytes_in_gb * 100) / 100;
}

export default {
	components: {
		"card-layout": CardLayout,
		"server-link": ServerLink
	},
	data: function() {
		return {
			hostingdomain_suffix: get_hostingdomain_suffix(),
			servers: [],
			node: null,
			disk_gb_total: 0,
			disk_gb_used: 0,
			disk_percent: 0,
			memory_gb_total: 0,
			memory_gb_used: 0,
			memory_percent: 0
		};
	},
	mounted: function() {
		this.update();
		this.handle = setInterval(() => this.update(), 5000);

		get_all_servers()
		.then(list => list.filter(s => s.user_node_id == this.$route.params.id))
		.then(list => this.servers = list);
	},
	beforeUnmount: function() {
		clearInterval(this.handle);
	},
	methods: {
		format_time: format_time,
		update: function() {
			get_by_id(this.$route.params.id)
			.then(n => this.node = n)
			.then(() => {
				this.disk_gb_total = get_gb_rounded(this.node.disk_size);
				this.disk_gb_used = get_gb_rounded(this.node.disk_used);
				this.disk_percent = parseInt(this.disk_gb_used / this.disk_gb_total * 100);
				this.memory_gb_total = get_gb_rounded(this.node.memory_size);
				this.memory_gb_used = get_gb_rounded(this.node.memory_used);
				this.memory_percent = parseInt(this.memory_gb_used / this.memory_gb_total * 100);
			});
		},
		save: function() {
			update_node(this.node)
			.then(() => this.update());
		}
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-server"></i> Node details
		</template>
		<h4>Details</h4>
		<table class="table" v-if="node">
			<tbody>
				<tr>
					<td>ID</td>
					<td>{{node.id}}</td>
				</tr>
				<tr>
					<td>Hostname</td>
					<td>{{node.name}}.{{hostingdomain_suffix}}</td>
				</tr>
				<tr>
					<td>IPv4</td>
					<td>{{node.ipv4}}</td>
				</tr>
				<tr>
					<td>IPv6</td>
					<td>{{node.ipv6}}</td>
				</tr>
				<tr>
					<td>State</td>
					<td>{{node.state}}</td>
				</tr>
				<tr v-if="servers.length == 0 && node.state == 'RUNNING'">
					<td>Actions</td>
					<td>
						<router-link class="btn btn-xs btn-danger" :to="'/nodes/' + node.id + '/delete'">
							<i class="fa fa-trash"></i>
							Delete
						</router-link>
					</td>
				</tr>
				<tr>
					<td>Alias</td>
					<td>
						<div class="btn-group w-100">
							<input type="text" class="form-control" v-model="node.alias"/>
							<button class="btn btn-xs btn-info" v-on:click="save()">
								<i class="fa fa-floppy-disk"></i>
							</button>
						</div>
					</td>
				</tr>
				<tr>
					<td>CPU Usage</td>
					<td>
						<div class="progress">
							<div class="progress-bar" v-bind:style="{ width: node.load_percent + '%' }">
								{{node.load_percent}}%
							</div>
						</div>
					</td>
				</tr>
				<tr>
					<td>Disk Usage</td>
					<td>
						{{disk_gb_used}}/{{disk_gb_total}} GB
						<div class="progress">
							<div class="progress-bar bg-warning" v-bind:style="{ width: disk_percent + '%' }">
								{{disk_percent}}%
							</div>
						</div>
					</td>
				</tr>
				<tr>
					<td>Memory Usage</td>
					<td>
						{{memory_gb_used}}/{{memory_gb_total}} GB
						<div class="progress">
							<div class="progress-bar bg-danger" v-bind:style="{ width: memory_percent + '%' }">
								{{memory_percent}}%
							</div>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
		<div v-if="this.node && this.node.state == 'RUNNING'">
			<h4>Servers</h4>
			<table class="table">
				<thead>
					<th>Name</th>
					<th>Created</th>
					<th>State</th>
				</thead>
				<tbody>
					<tr v-for="server in servers">
						<td>
							<server-link :server="server"/>
						</td>
						<td>
							{{format_time(server.created)}}
						</td>
						<td>
							{{server.state}}
						</td>
					</tr>
				</tbody>
			</table>
			<router-link class="btn btn-success" to="/mtservers/create">
				<i class="fa fa-plus"></i>
				Create server
			</router-link>
		</div>
	</card-layout>
	`
};
