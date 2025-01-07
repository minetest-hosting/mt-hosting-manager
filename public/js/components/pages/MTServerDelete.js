import CardLayout from "../layouts/CardLayout.js";
import { get_by_id, remove } from "../../api/mtserver.js";

export default {
	props: ["id"],
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			server: null,
			confirm_name: "",
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "list", name: "Servers", link: "/mtservers"
			},{
				icon: "list", name: "Server detail", link: `/mtservers/${this.id}`
			},{
				icon: "trash", name: "Server detail", link: `/mtservers/${this.id}/delete`
			}]
		};
	},
	mounted: async function() {
		this.server = await get_by_id(this.id);
	},
	methods: {
		remove: async function() {
			await remove(this.server);
			this.$router.push(`/nodes/${this.server.user_node_id}`);
		}
	},
	template: /*html*/`
	<card-layout title="Confirm server deletion" icon="trash">
		<table class="table" v-if="server">
			<tbody>
				<tr>
					<td>ID</td>
					<td>{{server.id}}</td>
				</tr>
				<tr>
					<td>Port</td>
					<td>{{server.port}}</td>
				</tr>
				<tr>
					<td>Name</td>
					<td>{{server.name}}</td>
				</tr>
				<tr>
					<td>Re-type name</td>
					<td>
						<input type="text" class="form-control" v-model="confirm_name"/>
					</td>
				</tr>
				<tr>
					<td>Delete</td>
					<td>
						<button class="btn btn-sm btn-danger" :disabled="confirm_name != server.name" v-on:click="remove">
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
