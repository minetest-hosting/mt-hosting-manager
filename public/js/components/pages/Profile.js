import CardLayout from "../layouts/CardLayout.js";
import { update } from "../../service/user.js";
import format_time from "../../util/format_time.js";
import { get_user_profile } from "../../service/user.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			profile: get_user_profile()
		};
	},
	methods: {
		format_time: format_time,
		save: function() {
			update(this.profile);
		}
	},
	template: /*html*/`
	<card-layout title="Profile" icon="user">
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
				<td>Balance</td>
				<td>
					&euro; {{profile.balance}}
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
