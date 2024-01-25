import CardLayout from "../layouts/CardLayout.js";
import NodeLink from "../NodeLink.js";
import ServerLink from "../ServerLink.js";
import NodeState from "../NodeState.js";
import ServerState from "../ServerState.js";
import NodeTypeSpec from "../NodeTypeSpec.js";

import { get_overview } from "../../api/overview.js";

import { get_user_id } from "../../service/login.js";
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
            nodes: []
		};
	},
    methods: {
        update: function() {
            get_overview(get_user_id())
            .then(od => {
                // populate nodetypes
                od.forEach(n => n.nodetype = get_nodetype(n.node_type_id));
                this.nodes = od;
            });
        }
    },
    mounted: function() {
        this.update();
        this.handle = setInterval(() => this.update(), 2000);
    },
	template: /*html*/`
	<card-layout title="Overview" icon="map" :breadcrumb="breadcrumb" :fullwidth="true" :flex="true">
        <div class="alert alert-info w-100" v-if="busy">
            <i class="fa fa-spinner fa-spin"></i>
			Loading overview data
        </div>
        <router-link class="btn btn-success" to="/nodes/create" v-if="nodes.length == 0 && !busy">
			<i class="fa fa-plus"></i>
			Create node
		</router-link>
        <div class="col-md-4" v-for="node in nodes" v-if="!busy">
            <div class="card" style="min-height: 250px;">
                <div class="card-header">
                    <node-link :node="node"/>
                    &nbsp;
                    <node-state :state="node.state"/>
                    &nbsp;
                    <node-type-spec :nodetype="node.nodetype"/>
                </div>
                <div class="card-body">
                    <ul>
                        <li v-for="server in node.servers">
                            <server-link :server="server"/>
                            &nbsp;
                            <server-state :state="server.state"/>
                        </li>
                        <li v-if="node.state == 'RUNNING'">
                            <router-link class="btn btn-sm btn-outline-success" :to="'/mtservers/create?node=' + node.id">
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
