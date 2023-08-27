
export default {
    props: ["id", "server"],
    computed: {
        name: function() {
            if (this.server) {
                return this.server.name;
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
