import CardLayout from "../layouts/CardLayout.js";
import NodeState from "../NodeState.js";
import ServerList from "../ServerList.js";

import { get_by_id, get_stats, update as update_node } from "../../api/node.js";
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
		"node-state": NodeState,
		"server-list": ServerList
	},
	data: function() {
		return {
			hostingdomain_suffix: get_hostingdomain_suffix(),
			servers: [],
			node: null,
			load_percent: 0,
			disk_gb_total: 0,
			disk_gb_used: 0,
			disk_percent: 0,
			memory_gb_total: 0,
			memory_gb_used: 0,
			memory_percent: 0
		};
	},
	mounted: function() {
		get_by_id(this.$route.params.id)
		.then(n => this.node = n);

		this.update_stats();
		this.handle = setInterval(() => this.update_stats(), 5000);

		get_all_servers()
		.then(list => list.filter(s => s.user_node_id == this.$route.params.id))
		.then(list => this.servers = list);
	},
	beforeUnmount: function() {
		clearInterval(this.handle);
	},
	methods: {
		format_time: format_time,
		update_stats: function() {
			get_stats(this.$route.params.id)
			.then(stats => {
				this.load_percent = stats.load_percent;
				this.disk_gb_total = get_gb_rounded(stats.disk_size);
				this.disk_gb_used = get_gb_rounded(stats.disk_used);
				this.disk_percent = parseInt(this.disk_gb_used / this.disk_gb_total * 100);
				this.memory_gb_total = get_gb_rounded(stats.memory_size);
				this.memory_gb_used = get_gb_rounded(stats.memory_used);
				this.memory_percent = parseInt(this.memory_gb_used / this.memory_gb_total * 100);
			});
		},
		save: function() {
			if (this.node && this.node.state == "RUNNING") {
				update_node(this.node);
			}
		}
	},
	template: /*html*/`
	<card-layout title="Node details" icon="server">
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
					<td>
						<node-state :state="node.state"/>
					</td>
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
							<div class="progress-bar" v-bind:style="{ width: load_percent + '%' }">
								{{load_percent}}%
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
			<server-list :list="servers"/>
			<router-link class="btn btn-success" to="/mtservers/create">
				<i class="fa fa-plus"></i>
				Create server
			</router-link>
		</div>
	</card-layout>
	`
};
