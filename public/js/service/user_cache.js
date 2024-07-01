import { search_user } from "../api/user.js";

// user_id -> user
const user_cache = {};
const user_cache_pending = {};

export const get_cached_user = user_id => {
    if (user_cache[user_id]) {
        return Promise.resolve(user_cache[user_id]);
    }
    if (user_cache_pending[user_id]) {
        return user_cache_pending[user_id];
    }

    const p = search_user({ user_id: user_id })
    .then(list => {
        delete user_cache_pending[user_id];

        if (list && list.length == 1) {
            user_cache[user_id] = list[0];
            return list[0];
        }
    });

    user_cache_pending[user_id] = p;
    return p;
};