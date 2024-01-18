import CardLayout from "../layouts/CardLayout.js";
import ServerLink from "../ServerLink.js";
import BackupState from "../BackupState.js";

import { get_by_id } from "../../api/backup_space.js";
import { get_by_backup_space_id, remove } from "../../api/backup.js";

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
                icon: "object-group", name: "Backup spaces", link: "/backup_spaces"
            },{
                icon: "object-group", name: `Backup space ${this.id}`, link: `/backup_spaces/${this.id}`
			}],
			space: null,
			backups: []
		};
	},
	methods: {
		format_time,
		format_size,
		update_backups: function() {
			get_by_backup_space_id(this.id).then(l => {
				l.sort((a,b) => a.created < b.created);
				this.backups = l;
			});
		},
		remove_backup: function(backup) {
			remove(backup).then(() => this.update_backups());
		}
	},
	mounted: function() {
		get_by_id(this.id).then(bs => this.space = bs);
		this.update_backups();
	},
	template: /*html*/`
	<card-layout title="Backup space detail" icon="object-group" :breadcrumb="breadcrumb">
		<table class="table table-condensed table-striped">
			<thead>
				<tr>
					<th>Created</th>
					<th>Server</th>
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
						<backup-state :state="backup.state"/>
						<span v-if="backup.state == 'COMPLETE'">
							&nbsp;
							{{format_size(backup.size)}}
						</span>
					</td>
					<td>
						<div class="btn-group">
							<button class="btn btn-sm btn-danger"
								v-on:click="remove_backup(backup)"
								:disabled="backup.state != 'COMPLETE' && backup.state != 'ERROR'">
								<i class="fa fa-trash"></i>
								Delete
							</button>
							<a class="btn btn-sm btn-secondary" :href="'/api/backup/' + backup.id">
								<i class="fa fa-download"></i>
								Download
							</a>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
	</card-layout>
	`
};
