import CardLayout from "../layouts/CardLayout.js";

import { get_all, create } from "../../api/backup_space.js";

export default {
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "object-group", name: "Backup spaces", link: "backup_spaces"
            }],
            spaces: [],
            new_name: "new backup space",
            busy: false
		};
	},
    mounted: async function() {
        this.busy = true;
        const s = await get_all();
        this.spaces = s;
        this.busy = false;
    },
    methods: {
        create: async function() {
            const bs = await create({ name: this.new_name });
            this.$router.push(`/backup_spaces/${bs.id}`);
        }
    },
	template: /*html*/`
	<card-layout title="Backup spaces" icon="object-group" :breadcrumb="breadcrumb">
        <table class="table table-condensed table-striped">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Retention</th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="space in spaces">
                    <td>
                        <router-link :to="'/backup_spaces/' + space.id">
                            {{space.name}}
                        </router-link>
                    </td>
                    <td>{{space.retention_days}} days</td>
                </tr>
                <tr>
                    <td>
                        <input class="form-control" type="text" placeholder="Name" v-model="new_name"/>
                    </td>
                    <td>
                        <a class="btn btn-success" v-if="!busy && new_name.length > 0" v-on:click="create">
                        <i class="fa fa-plus"></i>
                            Create backup space
                        </a>
                    </td>
                </tr>
            </tbody>
        </table>
        
	</card-layout>
	`
};
