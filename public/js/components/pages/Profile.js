import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";

import { update } from "../../service/user.js";
import format_time from "../../util/format_time.js";
import { get_user_profile } from "../../service/user.js";
import { get_currency_list } from "../../service/currency.js";

export default {
	components: {
		"card-layout": CardLayout,
		"currency-display": CurrencyDisplay
	},
	data: function() {
		return {
			profile: get_user_profile(),
			breadcrumb: [{
				icon: "home", name: "Home", link: "/"
			},{
				icon: "user", name: "Profile", link: "/profile"
			}]
		};
	},
	methods: {
		get_currency_list,
		format_time,
		save: function() {
			update(this.profile);
		}
	},
	template: /*html*/`
	<card-layout title="Profile" icon="user" :breadcrumb="breadcrumb">
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
					<currency-display :eurocents="profile.balance" :enable_warning="true"/>
				</td>
			</tr>
			<tr>
				<td>Warning balance threshold</td>
				<td>
					<input type="number" class="form-control" step="0.01" :value="profile.warn_balance / 100" v-on:input="e => profile.warn_balance = e.target.value * 100"/>
				</td>
			</tr>
			<tr>
				<td>Send balance warning mail</td>
				<td>
					<input type="checkbox" class="form-check-input" v-model="profile.warn_enabled"/>
				</td>
			</tr>
			<tr>
				<td>Preferred currency</td>
				<td>
					<select class="form-control" v-model="profile.currency">
						<option v-for="c in get_currency_list()" :key="c.id" :value="c.id">{{c.name}}</option>
					</select>
				</td>
			</tr>
			<tr>
				<td>Action</td>
				<td>
					<a class="btn btn-success w-100" v-on:click="save">
						<i class="fa-solid fa-floppy-disk"></i>
						Save
					</a>
				</td>
			</tr>
		</table>
	</card-layout>
    `
};
