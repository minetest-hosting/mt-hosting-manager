import CardLayout from "../layouts/CardLayout.js";
import NodeLink from "../NodeLink.js";
import ServerLink from "../ServerLink.js";
import { get_all, retry, remove } from "../../api/job.js";

export default {
	components: {
		"card-layout": CardLayout,
		"node-link": NodeLink,
		"server-link": ServerLink
	},
	data: function() {
		return {
			jobs: [],
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "play", name: "Jobs", link: "/jobs"
			}],
			handle: null
		};
	},
	mounted: function() {
		this.handle = setInterval(() => this.update(), 2000);
		this.update();
	},
	beforeUnmount: function() {
		clearInterval(this.handle);
	},
	methods: {
		update: function() {
			get_all().then(j => this.jobs = j);
		},
		retry: function(job) {
			retry(job)
			.then(() => this.update());
		},
		remove: function(job) {
			remove(job)
			.then(() => this.update());
		}
	},
	template: /*html*/`
	<card-layout title="Jobs" icon="play" :breadcrumb="breadcrumb">
		<table class="table table-condensed">
			<thead>
				<tr>
					<th>ID</th>
					<th>Type</th>
					<th>State</th>
					<th>Links</th>
					<th>Message</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="job in jobs">
					<td>{{job.id}}</td>
					<td>{{job.type}}</td>
					<td>{{job.state}}</td>
					<td>
						<node-link :id="job.user_node_id" v-if="job.user_node_id"/>
						<br>
						<server-link :id="job.minetest_server_id" v-if="job.minetest_server_id"/>
					</td>
					<td>{{job.message}}</td>
					<td>
						<div class="btn-group">
							<a class="btn btn-sm btn-outline-primary" v-on:click="retry(job)">
								<i class="fa-solid fa-rotate-right"></i>
								Retry
							</a>
							<a class="btn btn-sm btn-outline-danger" v-on:click="remove(job)">
								<i class="fa-solid fa-trash"></i>
								Remove
							</a>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
	</card-layout>
	`
};
