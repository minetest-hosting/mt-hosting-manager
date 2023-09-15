import CardLayout from "../layouts/CardLayout.js";
import { get_by_id, remove } from "../../api/node.js";

export default {
	props: ["id"],
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			node: null,
			confirm_alias: "",
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "server", name: "Nodes", link: "/nodes"
			},{
				icon: "server", name: "Node detail", link: `/nodes/${this.id}`
			},{
				icon: "trash", name: "Delete", link: `/nodes/${this.id}/delete`
			}]
		};
	},
	mounted: function() {
		get_by_id(this.id)
		.then(n => this.node = n);
	},
	methods: {
		remove: function() {
			remove(this.node)
			.then(() => this.$router.push("/nodes"));
		}
	},
	template: /*html*/`
	<card-layout title="Confirm node deletion" icon="trash" :breadcrumb="breadcrumb">
		<table class="table" v-if="node">
			<tbody>
				<tr>
					<td>ID</td>
					<td>{{node.id}}</td>
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
					<td>Alias</td>
					<td>{{node.alias}}</td>
				</tr>
				<tr>
					<td>Re-type alias</td>
					<td>
						<input type="text" class="form-control" v-model="confirm_alias"/>
					</td>
				</tr>
				<tr>
					<td>Delete</td>
					<td>
						<button class="btn btn-danger" :disabled="confirm_alias != node.alias" v-on:click="remove">
							<i class="fa fa-trash"></i>
							Delete
						</button>
					</td>
				</tr>
			</tbody>
		</table>
	</card-layout>
	`
};
