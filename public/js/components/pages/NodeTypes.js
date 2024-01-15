import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";
import NodeTypeState from "../NodeTypeState.js";
import NodeTypeSpec from "../NodeTypeSpec.js";

import { get_nodetypes, fetch_nodetypes } from "../../service/nodetype.js";
import { add } from "../../api/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout,
		"currency-display": CurrencyDisplay,
		"nodetype-state": NodeTypeState,
		"node-type-spec": NodeTypeSpec
	},
	data: function() {
		return {
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "server", name: "Nodetypes", link: "/node_types"
			}]
		};
	},
	computed: {
		nodetypes: get_nodetypes
	},
	methods: {
		on_import: function() {
			this.results = [];
			const files = this.$refs.upload.files;
			const file = files[0];
			file.text()
			.then(json => JSON.parse(json))
			.then(list => Promise.all(list.map(nt => add(nt))))
			.finally(() => fetch_nodetypes());
		}
	},
	template: /*html*/`
	<card-layout title="Nodetype" icon="server" :breadcrumb="breadcrumb" :fullwidth="true">
		<table class="table">
			<thead>
				<th>ID</th>
				<th>OrderID</th>
				<th>Name</th>
				<th>Location</th>
				<th>Specs</th>
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
					<td>
						<node-type-spec :nodetype="nt"/>
					</td>
					<td>{{nt.cpu_count}}</td>
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
		<div class="row">
			<div class="col-md-6">
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
			</div>
			<div class="col-md-6">
				<input ref="upload" type="file" class="form-control" v-on:change="on_import" accept=".json"/>
			</div>
		</div>
	</card-layout>
	`
};
