
export default {
    props: ["id", "node"],
    computed: {
        name: function() {
            if (this.node) {
                return `${this.node.alias} (${this.node.name})`;
            } else {
                return this.id;
            }
        },
        enabled: function() {
            // enabled if either just the id available or the node is running
            return (this.id || (this.node && this.node.state == "RUNNING"));
        }
    },
    template: /*html*/`
    <router-link :to="'/nodes/' + (id ? id : node.id)" v-if="enabled">
        <i class="fa fa-server"></i>
        {{name}}
    </router-link>
    <span v-else>
        {{name}}
    </span>
	`
};
