import ServerLink from "./ServerLink.js";
import NodeLink from "./NodeLink.js";
import ServerState from "./ServerState.js";
import format_time from "../util/format_time.js";

export default {
    props: ["list", "show_parent"],
    components: {
        "server-link": ServerLink,
        "node-link": NodeLink,
        "server-state": ServerState
    },
    methods: {
        format_time: format_time
    },
    template: /*html*/`
    <table class="table">
        <thead>
            <th>Name</th>
            <th v-if="show_parent">Parent node</th>
            <th>Created</th>
            <th>State</th>
        </thead>
        <tbody>
            <tr v-for="server in list" :key="server.id">
                <td>
                    <server-link :server="server"/>
                </td>
                <td v-if="show_parent">
                    <node-link :id="server.user_node_id"/>
                </td>
                <td>
                    {{format_time(server.created)}}
                </td>
                <td>
                    <server-state :state="server.state"/>
                </td>
            </tr>
        </tbody>
    </table>
    `
};
