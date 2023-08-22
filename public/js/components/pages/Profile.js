import CardLayout from "../layouts/CardLayout.js";
import { get_profile } from "../../api/user.js";
import format_time from "../../util/format_time.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			mail: null
		};
	},
	mounted: function() {
		get_profile().then(p => Object.assign(this, p));
	},
	methods: {
		format_time: format_time,
	},
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-user"></i> Profile
		</template>
		<table class="table" v-if="mail">
			<tr>
				<td>Mail</td>
				<td>{{mail}}</td>
			</tr>
			<tr>
				<td>Role</td>
				<td>{{role}}</td>
			</tr>
			<tr>
				<td>Type</td>
				<td>{{type}}</td>
			</tr>
			<tr>
				<td>Created</td>
				<td>
					{{format_time(created)}}
				</td>
			</tr>
		</table>
	</card-layout>
    `
};
