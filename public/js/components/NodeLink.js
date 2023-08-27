
export default {
    props: ["id", "node"],
    computed: {
        name: function() {
            if (this.node) {
                return `${this.node.alias} (${this.node.name})`;
            } else {
                return this.id;
            }
        }
    },
    template: /*html*/`
    <router-link :to="'/nodes/' + (id ? id : node.id)">
        <i class="fa fa-server"></i>
        {{name}}
    </router-link>
	`
};
