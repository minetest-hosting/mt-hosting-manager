import CardLayout from "../layouts/CardLayout.js";
import NodeLink from "../NodeLink.js";
import ServerLink from "../ServerLink.js";
import JobState from "../JobState.js";
import { get_all, retry, remove } from "../../api/job.js";
import format_time from '../../util/format_time.js';

export default {
	components: {
		"card-layout": CardLayout,
		"node-link": NodeLink,
		"server-link": ServerLink,
		"job-state": JobState
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
		format_time: format_time,
		update: async function() {
			this.jobs = await get_all();
		},
		retry: async function(job) {
			await retry(job);
			await this.update();
		},
		remove: async function(job) {
			await remove(job);
			await this.update();
		},
		rowClass: function(job) {
			const cl = {};
			switch (job.state) {
				case 'DONE_SUCCESS':
					cl["table-success"] = true;
					break;
				case 'DONE_FAILURE':
					cl["table-danger"] = true;
					break;
			}
			return cl;
		}
	},
	template: /*html*/`
	<card-layout title="Jobs" icon="play" :breadcrumb="breadcrumb" :fullwidth="true">
		<table class="table table-condensed">
			<thead>
				<tr>
					<th>ID</th>
					<th>Type</th>
					<th>Created</th>
					<th>Next run</th>
					<th>State</th>
					<th>Step</th>
					<th>Links</th>
					<th>Progress</th>
					<th>Message</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="job in jobs" v-bind:class="rowClass(job)">
					<td>{{job.id}}</td>
					<td>{{job.type}}</td>
					<td>{{format_time(job.created)}}</td>
					<td>{{format_time(job.next_run)}}</td>
					<td>
						<job-state :state="job.state"/>
					</td>
					<td>
						{{job.step}}
					</td>
					<td>
						<node-link :id="job.user_node_id" v-if="job.user_node_id"/>
						<br>
						<server-link :id="job.minetest_server_id" v-if="job.minetest_server_id"/>
					</td>
					<td>
						{{job.progress_percent}}
					</td>
					<td>{{job.message}}</td>
					<td>
						<div class="btn-group">
							<a class="btn btn-sm btn-outline-primary" v-on:click="retry(job)" v-if="job.state == 'DONE_FAILURE'">
								<i class="fa-solid fa-rotate-right"></i>
								Retry
							</a>
							<a class="btn btn-sm btn-outline-danger" v-on:click="remove(job)" v-if="job.state != 'RUNNING'">
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
