import CardLayout from "../layouts/CardLayout.js";
import { check, get_by_id } from "../../api/transaction.js";
import format_time from "../../util/format_time.js";

export default {
	components: {
		"card-layout": CardLayout
	},
    data: function() {
        return {
            transaction: null
        };
    },
    mounted: function() {
        get_by_id(this.$route.params.id)
        .then(tx => this.transaction = tx)
        .then(() => {
            check(this.transaction)
            .then(tx => this.transaction = tx)
            .then(() => fetch_profile());
        });
    },
    methods: {
        format_time: format_time,
    },
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-money-bill"></i> Finance details
		</template>
        <table class="table table-condensed" v-if="transaction">
            <tr>
                <td>State</td>
                <td>{{transaction.state}}</td>
            </tr>
        </table>
	</card-layout>
	`
};
