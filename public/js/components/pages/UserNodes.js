import CardLayout from "../layouts/CardLayout.js";
import { get_all } from "../../api/node.js";
import format_time from "../../util/format_time.js";
import { get_nodetype } from "../../service/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			nodes: []
		};
	},
	mounted: function() {
		this.update()
	},
	methods: {
		format_time: format_time,
		update: function() {
			get_all().then(nodes => this.nodes = nodes);
		},
		get_nodetype: get_nodetype
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-server"></i> Nodes
		</template>
		<table class="table">
			<thead>
				<th>Name</th>
				<th>State</th>
				<th>Created</th>
				<th>Node-Type</th>
			</thead>
			<tbody>
				<tr v-for="node in nodes">
					<td>
						<router-link :to="'/nodes/' + node.id">
							<i class="fa fa-server"></i>
							{{node.alias}} ({{node.name}})
						</router-link>
					</td>
					<td>{{node.state}}</td>
					<td>{{format_time(node.created)}}</td>
					<td>
						{{get_nodetype(node.nodetype_id)}}
					</td>
				</tr>
			</tbody>
		</table>
	</card-layout>
	`
};
