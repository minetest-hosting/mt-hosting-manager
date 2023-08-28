
export default {
    props: ["id", "server"],
    computed: {
        name: function() {
            if (this.server) {
                return this.server.name;
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
