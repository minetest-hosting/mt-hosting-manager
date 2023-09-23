import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";

import { get_all, create } from "../../api/transaction.js";
import format_time from "../../util/format_time.js";
import { get_user_profile } from "../../service/user.js";
import { get_max_balance } from "../../service/info.js";
import { get_balance } from "../../service/user.js";

export default {
	components: {
		"card-layout": CardLayout,
        "currency-display": CurrencyDisplay
	},
    data: function() {
        return {
            amount: 5,
            transactions: [],
            user: get_user_profile(),
            breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "money-bill", name: "Finance", link: "/finance"
            }]
        };
    },
    mounted: function() {
        this.update_payments();
    },
    methods: {
        format_time: format_time,
        get_max_balance: get_max_balance,
        new_payment: function() {
            create({ amount: Math.round(this.amount*100) })
            .then(r => window.location = r.url);
        },
        update_payments: function() {
            get_all().then(p => this.transactions = p);
        }
    },
    computed: {
        amount_sum_valid: function() {
            return get_balance() + (this.amount*100) <= get_max_balance();
        },
        amount_valid: function() {
            return this.amount_sum_valid && this.amount > 0;
        }
    },
	template: /*html*/`
	<card-layout title="Finance" icon="money-bill" :breadcrumb="breadcrumb">
        <h4>Current balance</h4>
        <table class="table table-condensed">
            <tr>
                <td>Balance</td>
                <td v-if="user">
                    <currency-display :eurocents="user.balance"/>
                </td>
            </tr>
            <tr>
                <td>Actions</td>
                <td>
                    <div class="input-group">
                        <span class="input-group-text" v-if="user">&euro;</span>
                        <input class="form-control" type="number" min="0" max="100" v-model="amount" v-bind:class="{'is-invalid':!amount_valid}"/>
                        <button class="btn btn-outline-primary" v-on:click="new_payment()" :disabled="!amount_valid">
                            <i class="fa-solid fa-plus"></i> Create new payment
                        </button>
                        <div class="invalid-feedback" v-if="!amount_sum_valid">
                            User-balance can't exceed <currency-display :eurocents="get_max_balance()"/>
                        </div>
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
                    <td>
                        <currency-display :eurocents="tx.amount"/>
                    </td>
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
