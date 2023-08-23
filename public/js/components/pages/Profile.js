import CardLayout from "../layouts/CardLayout.js";
import { get_profile, update_profile } from "../../api/user.js";
import format_time from "../../util/format_time.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			profile: null
		};
	},
	mounted: function() {
		get_profile().then(p => this.profile = p);
	},
	methods: {
		format_time: format_time,
		save: function() {
			update_profile(this.profile);
		}
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-user"></i> Profile
		</template>
		<table class="table" v-if="profile">
			<tr>
				<td>Mail</td>
				<td>{{profile.mail}}</td>
			</tr>
			<tr>
				<td>Role</td>
				<td>{{profile.role}}</td>
			</tr>
			<tr>
				<td>Type</td>
				<td>{{profile.type}}</td>
			</tr>
			<tr>
				<td>Created</td>
				<td>
					{{format_time(profile.created)}}
				</td>
			</tr>
			<tr>
				<td>Currency</td>
				<td>
					<select class="form-control" v-model="profile.currency">
						<option value="EUR">Euro</option>
						<option value="USD">US Dollar</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>Action</td>
				<td>
					<a class="btn btn-success" v-on:click="save">
						<i class="fa-solid fa-floppy-disk"></i>
						Save
					</a>
				</td>
			</tr>
		</table>
	</card-layout>
    `
};
