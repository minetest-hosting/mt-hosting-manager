import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";

import { get_nodetypes } from "../../service/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout,
		"currency-display": CurrencyDisplay
	},
	data: function() {
		return {
			nodetypes: get_nodetypes(),
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "server", name: "Nodetypes", link: "/node_types"
			}]
		};
	},
	template: /*html*/`
	<card-layout title="Nodetype" icon="server" :breadcrumb="breadcrumb">
		<table class="table">
			<thead>
				<th>ID</th>
				<th>OrderID</th>
				<th>Name</th>
				<th>Daily cost</th>
				<th>Provider</th>
				<th>Server-Type</th>
				<th>Actions</th>
			</thead>
			<tbody>
				<tr v-for="nt in nodetypes">
					<td>
						{{nt.id}}
						<span v-if="nt.state == 'ACTIVE'" class="badge text-bg-success">Active</span>
						<span v-if="nt.state == 'INACTIVE'" class="badge text-bg-info">Inactive</span>
						<span v-if="nt.state == 'DEPRECATED'" class="badge text-bg-warning">Deprecated</span>
					</td>
					<td>{{nt.order_id}}</td>
					<td>{{nt.name}}</td>
					<td>
						<currency-display :eurocents="nt.daily_cost"/>
					</td>
					<td>{{nt.provider}}</td>
					<td>{{nt.server_type}}</td>
					<td>
						<router-link class="btn btn-primary" :to="'/node_types/' + nt.id">
							<i class="fa fa-pen-to-square"></i>
							Edit
						</router-link>
					</td>
				</tr>
			</tbody>
		</table>
		<router-link class="btn btn-success" to="/node_types/new">
			<i class="fa fa-plus"></i>
			Create node-type
		</router-link>
	</card-layout>
	`
};
