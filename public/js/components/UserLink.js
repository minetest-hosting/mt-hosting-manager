import { get_cached_user } from "../service/user_cache.js";

export default {
    props: ["id"],
    data: function() {
        return {
            user: null
        };
    },
    mounted: function() {
        get_cached_user(this.id).then(u => this.user = u);
    },
    template: /*html*/`
        <span v-if="user">
            <router-link :to="'/users/' + user.id">
                {{user.name}}
            </router-link>
        </span>
        <span v-else>
            <router-link :to="'/users/' + id">
                {{id}}
            </router-link>
        </span>
    `
};