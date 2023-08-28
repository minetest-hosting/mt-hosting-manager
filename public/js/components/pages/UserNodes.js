import CardLayout from "../layouts/CardLayout.js";
import { get_all } from "../../api/node.js";
import format_time from "../../util/format_time.js";
import { get_nodetype } from "../../service/nodetype.js";
import NodeLink from "../NodeLink.js";

export default {
	components: {
		"card-layout": CardLayout,
		"node-link": NodeLink
	},
	data: function() {
		return {
			nodes: [],
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "server", name: "Nodes", link: "/nodes"
			}]
		};
	},
	mounted: function() {
		this.update();
	},
	methods: {
		format_time: format_time,
		update: function() {
			get_all().then(nodes => this.nodes = nodes);
		},
		get_nodetype: get_nodetype
	},
	template: /*html*/`
	<card-layout title="Nodes" icon="server" :breadcrumb="breadcrumb">
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
						<node-link :node="node"/>
					</td>
					<td>{{node.state}}</td>
					<td>{{format_time(node.created)}}</td>
					<td>
						{{get_nodetype(node.node_type_id).name}}
					</td>
				</tr>
			</tbody>
		</table>
		<router-link class="btn btn-success" to="/nodes/create">
			<i class="fa fa-plus"></i>
			Create node
		</router-link>
	</card-layout>
	`
};
