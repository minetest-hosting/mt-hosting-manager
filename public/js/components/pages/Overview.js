import CardLayout from "../layouts/CardLayout.js";
import NodeLink from "../NodeLink.js";
import ServerLink from "../ServerLink.js";
import NodeState from "../NodeState.js";
import ServerState from "../ServerState.js";
import NodeTypeSpec from "../NodeTypeSpec.js";

import { get_all, get_mtservers_by_nodeid } from "../../api/node.js";
import { get_nodetype } from "../../service/nodetype.js";

export default {
	components: {
		"card-layout": CardLayout,
        "node-link": NodeLink,
        "server-link": ServerLink,
        "node-state": NodeState,
        "server-state": ServerState,
        "node-type-spec": NodeTypeSpec
	},
	data: function() {
		return {
			breadcrumb: [{
                icon: "home", name: "Home", link: "/"
            },{
                icon: "map", name: "Overview", link: "/overview"
            }],
            busy: false,
            entries: []
		};
	},
    mounted: function() {
        get_all().then(nodes => {
            this.busy = true;
            const entries = [];
            const promises = nodes
            .filter(node => node.state == "RUNNING")
            .map(node => {
                const entry = {
                    node: node,
                    nodetype: get_nodetype(node.node_type_id),
                    servers: []
                };
                return get_mtservers_by_nodeid(node.id)
                .then(servers => {
                    entry.servers = servers
                        .filter(server => server.state == "RUNNING")
                        .sort((a,b) => a.id < b.id);
                    entries.push(entry);
                });
            });
            Promise.all(promises).then(() => {
                entries.sort((a,b) => a.node.id < b.node.id);
                this.entries = entries;
                this.busy = false;
            });
        });
    },
	template: /*html*/`
	<card-layout title="Overview" icon="map" :breadcrumb="breadcrumb" :fullwidth="true" :flex="true">
        <div class="alert alert-info w-100" v-if="busy">
            <i class="fa fa-spinner fa-spin"></i>
			Loading overview data
        </div>
        <router-link class="btn btn-success" to="/nodes/create" v-if="entries.length == 0 && !busy">
			<i class="fa fa-plus"></i>
			Create node
		</router-link>
        <div class="col-md-4" v-for="entry in entries" v-if="!busy">
            <div class="card" style="min-height: 250px;">
                <div class="card-header">
                    <node-link :node="entry.node"/>
                    &nbsp;
                    <node-state :state="entry.node.state"/>
                    &nbsp;
                    <node-type-spec :nodetype="entry.nodetype"/>
                </div>
                <div class="card-body">
                    <ul>
                        <li v-for="server in entry.servers">
                            <server-link :server="server"/>
                            &nbsp;
                            <server-state :state="server.state"/>
                        </li>
                        <li>
                            <router-link class="btn btn-sm btn-outline-success" :to="'/mtservers/create?node=' + entry.node.id">
                                <i class="fa fa-plus"></i>
                                Create server
                            </router-link>
                        </li>
                    </ul>
                </div>
            </div>
        </div>
	</card-layout>
	`
};
