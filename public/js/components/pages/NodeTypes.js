import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";
import NodeTypeState from "../NodeTypeState.js";

import { get_nodetypes } from "../../service/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout,
		"currency-display": CurrencyDisplay,
		"nodetype-state": NodeTypeState
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
	<card-layout title="Nodetype" icon="server" :breadcrumb="breadcrumb" :fullwidth="true">
		<table class="table">
			<thead>
				<th>ID</th>
				<th>OrderID</th>
				<th>Name</th>
				<th>Location</th>
				<th>CPU Count</th>
				<th>RAM GB</th>
				<th>Disk GB</th>
				<th>Daily cost</th>
				<th>Provider</th>
				<th>Server-Type</th>
				<th>Actions</th>
			</thead>
			<tbody>
				<tr v-for="nt in nodetypes">
					<td>
						{{nt.id}}
						<nodetype-state :state="nt.state"/>
					</td>
					<td>{{nt.order_id}}</td>
					<td>{{nt.name}}</td>
					<td>{{nt.location}} / {{nt.location_readable}}</td>
					<td>{{nt.cpu_count}}</td>
					<td>{{nt.ram_gb}}</td>
					<td>{{nt.disk_gb}}</td>
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
		<div class="btn-group">
			<router-link class="btn btn-success" to="/node_types/new">
				<i class="fa fa-plus"></i>
				Create node-type
			</router-link>
			<a class="btn btn-outline-secondary" href="api/nodetype?export=true">
				<i class="fa fa-download"></i>
				Export as json
			</a>
		</div>
	</card-layout>
	`
};
