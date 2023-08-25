import CardLayout from "../layouts/CardLayout.js";
import { get_all, create } from "../../api/transaction.js";
import format_time from "../../util/format_time.js";
import { get_user_profile } from "../../service/user.js";

export default {
	components: {
		"card-layout": CardLayout
	},
    data: function() {
        return {
            amount: 5,
            transactions: [],
            user: get_user_profile()
        };
    },
    mounted: function() {
        this.update_payments();
    },
    methods: {
        format_time: format_time,
        new_payment: function() {
            create({ amount: ""+this.amount })
            .then(r => window.location = r.url);
        },
        update_payments: function() {
            get_all().then(p => this.transactions = p);
        }
    },
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-money-bill"></i> Finance
		</template>
        <h4>Current balance</h4>
        <table class="table table-condensed">
            <tr>
                <td>Balance</td>
                <td v-if="user">
                    &euro; {{user.balance}}
                </td>
            </tr>
            <tr>
                <td>Actions</td>
                <td>
                    <div class="input-group">
                        <span class="input-group-text" v-if="user">&euro;</span>
                        <input class="form-control" type="number" min="0" max="100" v-model="amount"/>
                        <a class="btn btn-outline-primary" v-on:click="new_payment()">
                            <i class="fa-solid fa-plus"></i> Create new payment
                        </a>
                    </div>
                </td>
            </tr>
        </table>
        <hr>
        <h4>Payments</h4>
        <table class="table table-condensed">
            <thead>
                <tr>
                    <th>Date</th>
                    <th>Amount</th>
                    <th>State</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="tx in transactions">
                    <td>
                        <router-link :to="'/finance/detail/'+tx.id">
                            {{format_time(tx.created)}}
                        </router-link>
                    </td>
                    <td>&euro; {{tx.amount}}</td>
                    <td>
                        {{tx.state}}
                        <span class="badge bg-warning" v-if="tx.amount_refunded != '0'">
                            <i class="fa-solid fa-recycle"></i>
                            Refunded
                        </span>
                    </td>
                </tr>
            </tbody>
        </table>
	</card-layout>
	`
};
