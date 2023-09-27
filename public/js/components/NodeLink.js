import { get_all } from "../api/node.js";

const store = Vue.reactive({
    nodes: []
});

export default {
    props: ["id", "node"],
    created: function() {
        if (!store.nodes.length) {
            get_all().then(list => store.nodes = list);
        }
    },
    computed: {
        name: function() {
            let node = this.node;
            if (!node) {
                node = store.nodes.find(n => n.id == this.id);
            }
            if (node) {
                return `${node.alias} (${node.name})`;
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
