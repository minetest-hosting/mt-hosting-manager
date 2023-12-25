import CardLayout from "../layouts/CardLayout.js";
import NodeLink from "../NodeLink.js";
import ServerLink from "../ServerLink.js";

import { get_all, get_mtservers_by_nodeid } from "../../api/node.js";

export default {
	components: {
		"card-layout": CardLayout,
        "node-link": NodeLink,
        "server-link": ServerLink
	},
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "map", name: "Overview", link: "/overview"
            }],
            entries: []
		};
	},
    mounted: function() {
        get_all().then(nodes => {
            const entries = [];
            const promises = nodes.map(node => {
                const entry = {
                    node: node,
                    servers: []
                };
                return get_mtservers_by_nodeid(node.id)
                .then(servers => {
                    entry.servers = servers;
                    entries.push(entry);
                });
            });
            Promise.all(promises).then(() => {
                entries.sort((a,b) => a.id < b.id);
                this.entries = entries;
            });
        });
    },
	template: /*html*/`
	<card-layout title="Overview" icon="map" :breadcrumb="breadcrumb" :fullwidth="true">
        <div class="col-md-4" v-for="entry in entries">
            <div class="card" style="min-height: 250px;">
                <div class="card-header">
                    <node-link :node="entry.node"/>
                </div>
                <div class="card-body">
                    <ul>
                        <li v-for="server in entry.servers">
                            <server-link :server="server"/>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
	</card-layout>
	`
};
