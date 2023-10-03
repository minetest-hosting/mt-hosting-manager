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
        }
    },
    template: /*html*/`
    <router-link :to="'/mtservers/' + (id ? id : server.id)">
        <i class="fa fa-list"></i>
        {{name}}
    </router-link>
	`
};
