import CardLayout from "../layouts/CardLayout.js";
import CurrencyDisplay from "../CurrencyDisplay.js";

import { create } from "../../api/node.js";
import { get_nodetype, get_nodetypes } from "../../service/nodetype.js";
import { get_balance, fetch_profile } from "../../service/user.js";
import random_name from "../../util/random_name.js";

export default {
	components: {
		"card-layout": CardLayout,
        "currency-display": CurrencyDisplay
	},
    data: function() {
        return {
            nodetype_id: get_nodetypes()[0].id,
            alias: random_name(),
            busy: false
        };
    },
    methods: {
        create_node: function() {
            create({
                node_type_id: this.nodetype_id,
                alias: this.alias
            })
            .then(node => {
                this.$router.push(`/nodes/${node.id}`);
                fetch_profile();
            });
        }
    },
    computed: {
        balance: get_balance,
        enough_funds: function() {
            return this.balance >= (10 * this.nodetype.daily_cost);
        },
        nodetype: function() {
            return get_nodetype(this.nodetype_id);
        },
        available_nodetypes: function() {
            return get_nodetypes().filter(nt => nt.state == "ACTIVE");
        }
    },
	template: /*html*/`
	<card-layout title="Create node" icon="plus">
        <div class="row">
            <div class="col-12">
                <select v-model="nodetype_id" class="form-control" :disabled="busy">
                    <option v-for="nt in available_nodetypes" :value="nt.id">{{nt.name}}</option>
                </select>
            </div>
        </div>
        <hr>
        <table class="table" v-if="nodetype">
            <tbody>
                <tr>
                    <td>Name</td>
                    <td>{{nodetype.name}}</td>
                </tr>
                <tr>
                    <td>Description</td>
                    <td>{{nodetype.description}}</td>
                </tr>
                <tr>
                    <td>Location</td>
                    <td>{{nodetype.location_readable}}</td>
                </tr>
                <tr>
                    <td>CPU Count</td>
                    <td>
                        {{nodetype.cpu_count}}
                        <span class="badge bg-success" v-if="nodetype.dedicated">Dedicated</span>
                    </td>
                </tr>
                <tr>
                    <td>RAM [GB]</td>
                    <td>{{nodetype.ram_gb}}</td>
                </tr>
                <tr>
                    <td>Disk [GB]</td>
                    <td>{{nodetype.disk_gb}}</td>
                </tr>
                <tr>
                    <td>Daily cost</td>
                    <td>
                        <currency-display :eurocents="nodetype.daily_cost"/>
                        (for 30 days: <currency-display :eurocents="nodetype.daily_cost * 30"/>)
                        <div class="alert alert-warning" v-if="!enough_funds">
                            <i class="fa-solid fa-triangle-exclamation"></i>
                            Not enough funds: make sure you have at least enough funds to support 10 days of runtime
                        </div>
                    </td>
                </tr>
            </tbody>
        </table>
        <hr>
        <div class="row">
            <div class="col-12">
                <div class="input-group">
                    <span class="input-group-text">Alias</span>
                    <input class="form-control" placeholder="friendly nodename" type="text" v-model="alias" :disabled="busy"/>
                    <button class="btn btn-success" v-on:click="create_node()" :disabled="busy || !alias || !enough_funds">
                        <i class="fa fa-plus"></i>
                        Create new node
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                        <span class="badge bg-danger" v-if="!enough_funds">
                            <i class="fa-solid fa-triangle-exclamation"></i>
                            Not enough funds
                        </span>
                    </button>
                </div>
            </div>
        </div>
	</card-layout>
	`
};
