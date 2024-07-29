import { get_user_by_id } from "../api/user.js";

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

    const p = get_user_by_id(user_id)
    .then(user => {
        delete user_cache_pending[user_id];
        user_cache[user_id] = user;
        return user;
    });

    user_cache_pending[user_id] = p;
    return p;
};