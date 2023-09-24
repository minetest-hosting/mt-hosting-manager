import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";

import { check, get_by_id, refund } from "../../api/transaction.js";
import format_time from "../../util/format_time.js";
import { fetch_profile } from "../../service/user.js";
import { get_balance } from "../../service/user.js";
import { get_refund_amount } from "../../service/finance.js";

export default {
    props: ["id"],
	components: {
		"card-layout": CardLayout,
        "currency-display": CurrencyDisplay
	},
    data: function() {
        return {
            transaction: null,
            breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "money-bill", name: "Finance", link: "/finance"
            },{
                icon: "money-bill", name: "Transaction detail", link: `/finance/detail/${this.id}`
            }]
        };
    },
    mounted: function() {
        this.update();
    },
    methods: {
        format_time: format_time,
        update: function() {
            get_by_id(this.id)
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
        },
        get_refund_amount: get_refund_amount
    },
    computed: {
        balance: get_balance
    },
	template: /*html*/`
	<card-layout title="Finance details" icon="money-bill" :breadcrumb="breadcrumb">
        <table class="table table-condensed" v-if="transaction">
            <tr>
                <td>State</td>
                <td>{{transaction.state}}</td>
            </tr>
            <tr>
                <td>Amount</td>
                <td>
                    <currency-display :eurocents="transaction.amount"/>
                </td>
            </tr>
            <tr>
                <td>Refunded amount</td>
                <td>
                    <currency-display :eurocents="transaction.amount_refunded"/>
                </td>
            </tr>
            <tr>
                <td>Actions</td>
                <td>
                    <button class="btn btn-warning" v-on:click="refund" :disabled="transaction.amount_refunded > 0 || balance <= 0">
                        <i class="fa-solid fa-recycle"></i>
                        Refund
                        <currency-display :eurocents="get_refund_amount(transaction)"/>
                    </button>
                </td>
            </tr>
        </table>
	</card-layout>
	`
};
