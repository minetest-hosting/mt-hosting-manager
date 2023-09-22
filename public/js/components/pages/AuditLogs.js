import CardLayout from "../layouts/CardLayout.js";
import NodeLink from "../NodeLink.js";
import ServerLink from "../ServerLink.js";
import PaymentLink from "../PaymentLink.js";
import CurrencyDisplay from "../CurrencyDisplay.js";

import { search_audit_logs } from "../../api/audit_log.js";
import format_time from "../../util/format_time.js";

const store = Vue.reactive({
    from: new Date(Date.now() - (3600*1000*2)),
    to: new Date(Date.now() + (3600*1000*1)),
	breadcrumb: [{
		icon: "home", name: "Home", link: "/"
	}, {
		icon: "rectangle-list", name: "Audit-Logs", link: "/audit-logs"
	}],
	list: [],
	busy: false
});

export default {
	components: {
		"card-layout": CardLayout,
		"node-link": NodeLink,
		"server-link": ServerLink,
		"payment-link": PaymentLink,
		"currency-display": CurrencyDisplay
	},
	data: () => store,
	methods: {
		format_time: format_time,
		search: function() {
			this.busy = true;
			search_audit_logs({
				from_timestamp: Math.floor(+this.from/1000),
				to_timestamp: Math.floor(+this.to/1000)
			})
			.then(l => this.list = l)
			.finally(() => this.busy = false);
		},
	},
	watch: {
		"from": "search",
		"to": "search"
	},
	created: function() {
		this.search();
	},
	template: /*html*/`
	<card-layout title="Audit-Logs" icon="rectangle-list" :breadcrumb="breadcrumb" fullwidth="true">
		<div class="row">
			<div class="col-5">
				<label>From</label>
				<vue-datepicker v-model="from"/>
			</div>
			<div class="col-5">
				<label>To</label>
				<vue-datepicker v-model="to"/>
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
		<div class="alert alert-warning" v-if="list.length >= 1000">
			<i class="fa-solid fa-triangle-exclamation"></i>
			<b>Warning:</b> more than 1000 results available and only 1000 fetched,
			reduce the time-window to get all relevant entries
		</div>
		<table class="table table-condensed table-striped">
			<thead>
				<tr>
					<th>Time</th>
					<th>Type</th>
					<th>Action</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="log in list" :key="log.id">
					<td>{{format_time(log.timestamp)}}</td>
					<td>{{log.type}}</td>
					<td>
						<node-link :id="log.user_node_id" v-if="log.user_node_id"/>
						<server-link :id="log.minetest_server_id" v-if="log.minetest_server_id"/>
						<payment-link :id="log.payment_transaction_id" v-if="log.payment_transaction_id"/>
						<span v-if="log.amount" class="badge bg-success">
							<currency-display :eurocents="log.amount"/>
						</span>
					</td>
				</tr>
			</tbody>
		</table>
	</card-layout>
	`
};
