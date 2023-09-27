import { get_all } from "../api/mtserver.js";

const store = Vue.reactive({
    servers: []
});

export default {
    props: ["id", "server"],
    created: function() {
        if (!store.servers.length) {
            get_all().then(l => store.servers = l);
        }
    },
    computed: {
        name: function() {
            let server = this.server;
            if (!server) {
                server = store.servers.find(s => s.id == this.id);
            }
            if (server) {
                return server.name;
            } else {
                return this.id;
            }
        },
        enabled: function() {
            // enabled if either just the id available or the server is running
            return (this.id || (this.server && this.server.state == "RUNNING"));
        }
    },
    template: /*html*/`
    <router-link :to="'/mtservers/' + (id ? id : server.id)" v-if="enabled">
        <i class="fa fa-list"></i>
        {{name}}
    </router-link>
    <span v-else>
        {{name}}
    </span>
	`
};
