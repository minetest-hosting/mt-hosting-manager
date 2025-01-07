import CardLayout from "../layouts/CardLayout.js";
import ServerLink from "../ServerLink.js";
import BackupState from "../BackupState.js";

import { get_all, remove } from "../../api/backup.js";

import format_time from "../../util/format_time.js";
import format_size from "../../util/format_size.js";

export default {
	props: ["id"],
	components: {
		"card-layout": CardLayout,
		"server-link": ServerLink,
		"backup-state": BackupState
	},
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "object-group", name: "Backups", link: "/backup"
			}],
			backups: []
		};
	},
	methods: {
		format_time,
		format_size,
		update_backups: async function() {
			const l = await get_all();
			l.sort((a,b) => a.created < b.created);
			this.backups = l;
		},
		remove_backup: async function(backup) {
			await remove(backup);
			this.update_backups();
		}
	},
	mounted: function() {
		this.update_backups();
	},
	template: /*html*/`
	<card-layout title="Backups" icon="object-group" :breadcrumb="breadcrumb">
		<table class="table table-condensed table-striped">
			<thead>
				<tr>
					<th>Created</th>
					<th>Server</th>
					<th>Size</th>
					<th>State</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="backup in backups">
					<td>{{format_time(backup.created)}}</td>
					<td>
						<server-link :id="backup.minetest_server_id"/>
					</td>
					<td>
						<span v-if="backup.state == 'COMPLETE'">
							{{format_size(backup.size)}}
						</span>
					</td>
					<td>
						<backup-state :state="backup.state"/>
					</td>
					<td>
						<div class="btn-group">
							<button class="btn btn-sm btn-danger"
								v-on:click="remove_backup(backup)"
								:disabled="backup.state != 'COMPLETE' && backup.state != 'ERROR'">
								<i class="fa fa-trash"></i>
								Delete
							</button>
							<a class="btn btn-sm btn-secondary" :href="'/api/backup/' + backup.id + '/download'">
								<i class="fa fa-download"></i>
								Download
							</a>
							<router-link class="btn btn-sm btn-outline-secondary" :to="'/mtservers/create?restore_from=' + backup.id">
								<i class="fa fa-play"></i>
								Restore
							</router-link>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
	</card-layout>
	`
};
