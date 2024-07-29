import format_time from "../../util/format_time.js";
import { get_user_by_id } from "../../api/user.js";

import CurrencyDisplay from "../CurrencyDisplay.js";
import CardLayout from "../layouts/CardLayout.js";

export default {
    props: ["id"],
	components: {
		"card-layout": CardLayout,
        "currency-display": CurrencyDisplay
	},
    methods: {
        format_time
    },
	data: function() {
		return {
			breadcrumb: [{
                icon: "users", name: "Users", link: `/users`
            },{
                icon: "user", name: this.id, link: `/users/${this.id}`
            }],
            user: null
		};
	},
    mounted: async function() {
        this.user = await get_user_by_id(this.id);
    },
	template: /*html*/`
	<card-layout :title="'User ' + this.user.name" icon="user" :breadcrumb="breadcrumb" v-if="user">
        <table class="table table-striped table-condensed">
            <tbody>
                <tr>
                    <td>Name</td>
                    <td>{{user.name}}</td>
                </tr>
                <tr>
                    <td>ID</td>
                    <td>{{user.id}}</td>
                </tr>
                <tr>
                    <td>State</td>
                    <td>{{user.state}}</td>
                </tr>
                <tr>
                    <td>Type</td>
                    <td>{{user.type}}</td>
                </tr>
                <tr>
                    <td>Role</td>
                    <td>{{user.role}}</td>
                </tr>
                <tr>
                    <td>Created</td>
                    <td>{{format_time(user.created)}}</td>
                </tr>
                <tr>
                    <td>Balance</td>
                    <td>
                        <currency-display :eurocents="user.balance"/>
                    </td>
                </tr>
            </tbody>
        </table>
	</card-layout>
	`
};
