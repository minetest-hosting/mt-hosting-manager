import CardLayout from "../layouts/CardLayout.js";
import store from "../../store/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout
	},
    data: function() {
        return {
            nodetype_id: store.nodetypes[0].id
        };
    },
    computed: {
        nodetype: function() {
            return store.nodetypes.find(nt => nt.id == this.nodetype_id);
        },
        available_nodetypes: function() {
            return store.nodetypes.filter(nt => nt.state == "ACTIVE");
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
                <select v-model="nodetype_id" class="form-control">
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
            <tr>
                <td></td>
                <td></td>
            </tr>
        </table>
        <hr>
        <div class="row">
            <div class="col-12">
                <a class="btn btn-success">
                    <i class="fa fa-plus"></i>
                    Create new node
                </a>
            </div>
        </div>
	</card-layout>
	`
};
