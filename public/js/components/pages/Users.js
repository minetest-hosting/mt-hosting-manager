import format_time from "../../util/format_time.js";
import { get_users } from "../../api/user.js";

import CurrencyDisplay from "../CurrencyDisplay.js";
import CardLayout from "../layouts/CardLayout.js";

export default {
	components: {
		"card-layout": CardLayout,
        "currency-display": CurrencyDisplay
	},
    methods: {
        format_time
    },
	data: function() {
		return {
			breadcrumb: [{icon: "users", name: "Users", link: "/users"}],
            users: []
		};
	},
    mounted: async function() {
        this.users = await get_users();
    },
	template: /*html*/`
	<card-layout title="Users" icon="users" :breadcrumb="breadcrumb">
        <table class="table table-striped table-condensed">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>ID</th>
                    <th>State</th>
                    <th>Type</th>
                    <th>Role</th>
                    <th>Created</th>
                    <th>Last login</th>
                    <th>Balance</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="user in users">
                    <td>
                        <router-link :to="'/users/' + user.id">
                            {{user.name}}
                        </router-link>
                    </td>
                    <td>{{user.id}}</td>
                    <td>{{user.state}}</td>
                    <td>{{user.type}}</td>
                    <td>{{user.role}}</td>
                    <td>{{format_time(user.created)}}</td>
                    <td>{{format_time(user.lastlogin)}}</td>
                    <td>
                        <currency-display :eurocents="user.balance"/>
                    </td>
                </tr>
            </tbody>
        </table>
	</card-layout>
	`
};
