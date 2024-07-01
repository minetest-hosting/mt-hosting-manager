import { get_cached_user } from "../service/user_cache.js";

export default {
    props: ["user_id"],
    data: function() {
        return {
            user: null
        };
    },
    mounted: function() {
        get_cached_user(this.user_id).then(u => this.user = u);
    },
    template: /*html*/`
        <span v-if="user">
            {{user.name}}
            <span class="badge bg-info">{{user.type}}</span>
        </span>
        <span v-else>{{user_id}}</span>
    `
};