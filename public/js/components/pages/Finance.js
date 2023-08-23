import CardLayout from "../layouts/CardLayout.js";
import { get_all, create } from "../../api/transaction.js";

export default {
	components: {
		"card-layout": CardLayout
	},
    data: function() {
        return {
            amount: 5,
            payment_url: "",
            payments: []
        };
    },
    mounted: function() {
        this.update_payments();
    },
    methods: {
        new_payment: function() {
            create({
                amount: ""+this.amount
            })
            .then(r => {
                this.payment_url = r.url;
            });
        },
        update_payments: function() {
            get_all().then(p => this.payments = p);
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
                <td>&euro; 10.0</td>
            </tr>
            <tr>
                <td>Actions</td>
                <td>
                    <div class="input-group">
                        <span class="input-group-text">$</span>
                        <input class="form-control" type="number" min="0" max="100" v-model="amount"/>
                        <a class="btn btn-outline-primary" v-on:click="new_payment()">
                            <i class="fa-solid fa-plus"></i> Create new payment
                        </a>
                    </div>
                    <a class="btn btn-success" v-if="payment_url" :href="payment_url">
                        <i class="fa-solid fa-cart-shopping"></i>
                        Open payment page
                    </a>
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
                <tr>
                    <td></td>
                    <td></td>
                    <td></td>
                </tr>
            </tbody>
        </table>
	</card-layout>
	`
};
