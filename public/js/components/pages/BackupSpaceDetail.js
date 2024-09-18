import CardLayout from "../layouts/CardLayout.js";
import ServerLink from "../ServerLink.js";
import BackupState from "../BackupState.js";

import { get_by_id, update, remove as remove_space } from "../../api/backup_space.js";
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
			space_save_ok: false,
			space_remove_confirm: "",
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
		},
		save_space: function() {
			update(this.space).then(() => this.space_save_ok = true);
		},
		remove_space: function() {
			remove_space(this.space)
			.then(() => this.$router.push("/backup_spaces"));
		}
	},
	mounted: function() {
		get_by_id(this.id).then(bs => this.space = bs);
		this.update_backups();
	},
	template: /*html*/`
	<card-layout title="Backup space detail" icon="object-group" :breadcrumb="breadcrumb">
		<table class="table" v-if="space">
			<tr>
				<td>Name</td>
				<td>
					<input class="form-control" type="text" placeholder="Name" v-model="space.name"/>
				</td>
			</tr>
			<tr>
				<td>Retention (days)</td>
				<td>
					<input class="form-control" type="number" min="7" placeholder="Retention (days)" v-model="space.retention_days"/>
				</td>
			</tr>
			<tr>
				<td>Save</td>
				<td>
					<a class="btn btn-success btn-sm w-100" v-on:click="save_space">
						<i class="fa fa-floppy-disk"></i>
						Save
						<i class="fa fa-check" color="green" v-if="space_save_ok"></i>
					</a>
				</td>
			</tr>
			<tr>
				<td>Delete space</td>
				<td class="input-group">
					<input class="form-control" type="text" placeholder="Re-type name to confirm" v-model="space_remove_confirm"/>
					<button class="input-group-addon btn btn-danger btn-sm" :disabled="space_remove_confirm != space.name" v-on:click="remove_space">
						<i class="fa fa-trash"></i>
						Delete backup-space
					</button>
				</td>
			</tr>
		</table>
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
