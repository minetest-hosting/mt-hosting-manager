import CardLayout from "../layouts/CardLayout.js";
import { get_by_id, add, update, remove } from "../../api/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			nt: null
		};
	},
	mounted: function() {
		const id = this.$route.params.id;
		if (id != "new") {
			get_by_id(this.$route.params.id).then(nt => this.nt = nt);
		} else {
			this.nt = {
				id: "",
				state: "INACTIVE",
				order_id: 0
			};
		}
	},
	methods: {
		remove: function() {
			remove(this.nt).then(() => this.$router.push("/node_types"));
		},
		save: function() {
			if (this.nt.id == "") {
				// create new
				add(this.nt).then(nt => Object.assign(this.nt, nt))
				.then(() => this.$router.push("/node_types"));
			} else {
				// update existing
				update(this.nt)
				.then(() => this.$router.push("/node_types"));
			}
		}
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-server"></i> Nodetype
		</template>
		<table class="table" v-if="nt">
			<tr>
				<td>ID</td>
				<td>
					<input type="text" readonly="true" disabled="true" class="form-control" v-model="nt.id"/>
				</td>
			</tr>
			<tr>
				<td>State</td>
				<td>
					<select name="state" class="form-control" v-model="nt.state">
						<option value="ACTIVE">Active</option>
						<option value="INACTIVE">Inactive</option>
						<option value="DEPRECATED">Deprecated</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>Order ID</td>
				<td>
					<input type="number" class="form-control" v-model="nt.order_id"/>
				</td>
			</tr>
			<tr>
				<td>Provider</td>
				<td>
					<select name="provider" class="form-control" v-model="nt.provider">
						<option value="HETZNER">Hetzner</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>Server Type</td>
				<td>
					<input type="text" class="form-control" v-model="nt.server_type"/>
				</td>
			</tr>
			<tr>
				<td>Location</td>
				<td>
					<input type="text" class="form-control" v-model="nt.location"/>
				</td>
			</tr>
			<tr>
				<td>Name</td>
				<td>
					<input type="text" class="form-control" v-model="nt.name"/>
				</td>
			</tr>
			<tr>
				<td>Description</td>
				<td>
					<textarea rows="8" class="form-control" v-model="nt.description"></textarea>
				</td>
			</tr>
			<tr>
				<td>Daily cost [&euro;]</td>
				<td>
					<input type="text" class="form-control" v-model="nt.daily_cost"/>
				</td>
			</tr>
			<tr>
				<td>Max recommended instances</td>
				<td>
					<input type="number" min="0" class="form-control" v-model="nt.max_recommended_instances"/>
				</td>
			</tr>
			<tr>
				<td>Max instances</td>
				<td>
					<input type="number" min="0" class="form-control" v-model="nt.max_instances"/>
				</td>
			</tr>
			<tr>
				<td>
					<div class="d-grid gap-2">
						<button class="btn btn-success" v-on:click="save">
							<i class="fa fa-floppy-disk"></i>
							Save
						</button>
					</div>    
				</td>
				<td>
					<div class="d-grid gap-2">
						<button class="btn btn-danger" v-on:click="remove">
							<i class="fa fa-floppy-disk"></i>
							Delete
						</button>
					</div>    
				</td>
			</tr>
		</table>
	</card-layout>
	`
};
