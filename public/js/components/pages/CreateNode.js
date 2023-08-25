import CardLayout from "../layouts/CardLayout.js";
import { create } from "../../api/node.js";
import { get_nodetype, get_nodetypes } from "../../service/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout
	},
    data: function() {
        return {
            nodetype_id: get_nodetypes()[0].id,
            alias: "",
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
            });
        }
    },
    computed: {
        nodetype: function() {
            return get_nodetype(this.nodetype_id);
        },
        available_nodetypes: function() {
            return get_nodetypes().filter(nt => nt.state == "ACTIVE");
        }
    },
	template: /*html*/`
	<card-layout>
		<template #title>
			<i class="fa fa-server"></i>
            <i class="fa fa-plus"></i>
            Create node
		</template>
        <div class="row">
            <div class="col-12">
                <select v-model="nodetype_id" class="form-control" :disabled="busy">
                    <option v-for="nt in available_nodetypes" :value="nt.id">{{nt.name}}</option>
                </select>
            </div>
        </div>
        <hr>
        <table class="table" v-if="nodetype">
            <tr>
                <td>Name</td>
                <td>{{nodetype.name}}</td>
            </tr>
            <tr>
                <td>Description</td>
                <td>{{nodetype.description}}</td>
            </tr>
            <tr>
                <td>Daily cost</td>
                <td>&euro; {{nodetype.daily_cost}}</td>
            </tr>
        </table>
        <hr>
        <div class="row">
            <div class="col-12">
                <div class="input-group">
                    <span class="input-group-text">Alias</span>
                    <input class="form-control" placeholder="friendly nodename" type="text" v-model="alias" :disabled="busy"/>
                    <button class="btn btn-success" v-on:click="create_node()" :disabled="busy || !alias">
                        <i class="fa fa-plus"></i>
                        Create new node
                        <i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
                    </button>
                </div>
            </div>
        </div>
	</card-layout>
	`
};
