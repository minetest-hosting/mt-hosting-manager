import CardLayout from "../layouts/CardLayout.js";

import { get_by_id } from "../../api/backup_space.js";
import { get_by_backup_space_id } from "../../api/backup.js";

export default {
	props: ["id"],
	components: {
		"card-layout": CardLayout
	},
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "object-group", name: "Backup spaces", link: "/backup_spaces"
            },{
                icon: "object-group", name: `Backup space ${this.id}`, link: `/backup_spaces/${this.id}`
			}],
			space: null,
			backups: []
		};
	},
	mounted: function() {
		get_by_id(this.id).then(bs => this.space = bs);
		get_by_backup_space_id(this.id).then(l => this.backups = l);
	},
	template: /*html*/`
	<card-layout title="Backup space detail" icon="object-group" :breadcrumb="breadcrumb">
	</card-layout>
	`
};
