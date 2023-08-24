import CardLayout from "../layouts/CardLayout.js";
import { check, get_by_id, refund } from "../../api/transaction.js";
import format_time from "../../util/format_time.js";
import { fetch_profile } from "../../service/user.js";

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
        this.update();
    },
    methods: {
        format_time: format_time,
        update: function() {
            get_by_id(this.$route.params.id)
            .then(tx => this.transaction = tx)
            .then(() => {
                if (this.transaction.state != "SUCCESS") {
                    return check(this.transaction)
                    .then(tx => this.transaction = tx)
                    .then(() => fetch_profile());
                }
            });
        },
        refund: function() {
            refund(this.transaction)
            .then(() => this.update())
            .then(() => fetch_profile());
        }
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
            <tr>
                <td>Amount</td>
                <td>&euro; {{transaction.amount}}</td>
            </tr>
            <tr>
                <td>Refunded amount</td>
                <td>&euro; {{transaction.amount_refunded}}</td>
            </tr>
            <tr>
                <td>Actions</td>
                <td>
                    <button class="btn btn-warning" v-on:click="refund" :disabled="transaction.amount_refunded != '0'">
                        <i class="fa-solid fa-recycle"></i>
                        Refund
                    </button>
                </td>
            </tr>
        </table>
	</card-layout>
	`
};
