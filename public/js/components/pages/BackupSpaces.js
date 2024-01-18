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
            busy: false
		};
	},
    mounted: function() {
        this.busy = true;
        get_all().then(s => {
            this.spaces = s;
            this.busy = false;
        });
    },
    methods: {
        create: function() {
            create({
                name: "New backup space"
            })
            .then(bs => this.$router.push(`/backup_spaces/${bs.id}`));
        }
    },
	template: /*html*/`
	<card-layout title="Backup spaces" icon="object-group" :breadcrumb="breadcrumb">
        <ul>
            <li v-for="space in spaces">
                <router-link :to="'/backup_spaces/' + space.id">
                    {{space.name}} ({{space.id}})
                </router-link>
            </li>
        </ul>
        <a class="btn btn-success" v-if="!busy && spaces.length == 0" v-on:click="create">
            <i class="fa fa-plus"></i>
            Create backup space
        </a>
	</card-layout>
	`
};
