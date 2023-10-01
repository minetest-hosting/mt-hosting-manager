import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";
import UserSearch from "../UserSearch.js";

import { search_transaction, create } from "../../api/transaction.js";
import format_time from "../../util/format_time.js";
import { get_max_balance } from "../../service/info.js";
import { get_balance } from "../../service/user.js";
import { has_role } from "../../service/login.js";

export default {
	components: {
		"card-layout": CardLayout,
        "currency-display": CurrencyDisplay,
        "user-search": UserSearch
	},
    data: function() {
        return {
            amount: 5,
            busy: false,
            transactions: [],
            user: null,
            from: new Date(Date.now() - (3600*1000*24*30)),
            to: new Date(Date.now() + (3600*1000*1)),        
            breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "money-bill", name: "Finance", link: "/finance"
            }]
        };
    },
    mounted: function() {
        this.search();
    },
    methods: {
        format_time: format_time,
        get_max_balance: get_max_balance,
        has_role: has_role,
        new_payment: function(type) {
            this.busy = true;
            create({
                amount: Math.round(this.amount*100),
                type: type
            })
            .then(ctx => window.location = ctx.payment_url);
        },
        search: function() {
            this.busy = true;
            search_transaction({
                from_timestamp: Math.floor(+this.from/1000),
				to_timestamp: Math.floor(+this.to/1000),
				user_id: this.user ? this.user.id : null
            })
            .then(p => this.transactions = p)
            .finally(() => this.busy = false);
        }
    },
    computed: {
        balance: get_balance,
        min_sum_error: function() {
            return this.amount < 5;
        },
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
            <thead>
                <tr>
                    <td>Balance</td>
                    <td>
                        <currency-display :eurocents="balance" :enable_warning="true"/>
                    </td>
                </tr>
                <tr>
                    <td>Payment</td>
                    <td>
                        <div class="input-group">
                            <span class="input-group-text">&euro;</span>
                            <input class="form-control" type="number" min="5" max="100" v-model="amount" v-bind:class="{'is-invalid':!amount_valid||min_sum_error}"/>
                            <button class="btn btn-outline-primary" v-on:click="new_payment('WALLEE')" :disabled="busy||!amount_valid||min_sum_error">
                                <i class="fa-brands fa-cc-visa"></i>
                                <i class="fa-brands fa-paypal"></i>
                                Pay
                            </button>
                            <button class="btn btn-outline-primary" v-on:click="new_payment('COINBASE')" :disabled="busy||!amount_valid||min_sum_error">
                                <i class="fa-brands fa-bitcoin"></i>
                                <i class="fa-brands fa-ethereum"></i>
                                Pay with crypto
                            </button>
                            <div class="invalid-feedback" v-if="!amount_sum_valid">
                                User-balance can't exceed <currency-display :eurocents="get_max_balance()"/>
                            </div>
                            <div class="invalid-feedback" v-if="min_sum_error">
                                Minimum payment: <currency-display eurocents="500"/>
                            </div>
                        </div>
                    </td>
                </tr>
            </thead>
        </table>
        <hr>
        <h4>Payment history</h4>
        <div class="row">
			<div class="col-4">
				<label>From</label>
				<vue-datepicker v-model="from"/>
			</div>
			<div class="col-4">
				<label>To</label>
				<vue-datepicker v-model="to"/>
			</div>
			<div class="col-2">
                <div v-if="has_role('ADMIN')">
                    <label>User</label>
                    <user-search v-model="user"/>
                </div>
			</div>
			<div class="col-2">
				<label>Search</label>
				<button class="btn btn-primary w-100" v-on:click="search">
					<i class="fa fa-magnifying-glass"></i>
					Search
					<i class="fa fa-spinner fa-spin" v-if="busy"></i>
				</button>
			</div>
		</div>
		<hr>
        <table class="table table-condensed">
            <thead>
                <tr>
                    <th>Date</th>
                    <th>Amount</th>
                    <th>Type</th>
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
                    <td>{{tx.type}}</td>
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
