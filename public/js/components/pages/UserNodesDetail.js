import CardLayout from "../layouts/CardLayout.js";
import { get_by_id } from "../../api/node.js";

const bytes_in_gb = 1000 * 1000 * 1000;

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			node: null,
			disk_gb_total: 0
		};
	},
	mounted: function() {
		this.update();
	},
	methods: {
		update: function() {
			get_by_id(this.$route.params.id)
			.then(n => this.node = n)
			.then(() => {
				this.disk_gb_total = parseInt(this.node.disk_size / bytes_in_gb);
			})
		}
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-server"></i> Node details
		</template>
		<table class="table table-condensed" v-if="node">
			<tr>
				<td>ID</td>
				<td>{{node.id}}</td>
			</tr>
			<tr>
				<td>State</td>
				<td>{{node.state}}</td>
			</tr>
			<tr>
				<td>Alias</td>
				<td>{{node.alias}}</td>
			</tr>
			<tr>
				<td>CPU Usage</td>
				<td>
					<div class="progress">
						<div class="progress-bar"
							v-bind:style="{ width: node.load_percent + '%' }">
							{{node.load_percent}}%
							</div>
					</div>
				</td>
			</tr>
			<tr>
				<td>Disk Usage</td>
				<td>
					0/{{disk_gb_total}} GB
					<div class="progress">
						<div class="progress-bar bg-warning" style="width: 10%">
							10%
						</div>
					</div>
				</td>
			</tr>
			<tr>
				<td>Memory Usage</td>
				<td>
					0.2/2.0 GB
					<div class="progress">
						<div class="progress-bar bg-danger" style="width: 10%">
							10%
						</div>
					</div>
				</td>
			</tr>

		</table>
	</card-layout>
	`
};
