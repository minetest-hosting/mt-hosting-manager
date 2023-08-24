
import { get_profile, update_profile } from "../api/user.js";
import store from '../store/user.js';
import events, { EVENT_LOGGED_IN } from '../events.js';

events.on(EVENT_LOGGED_IN, function() {
    get_profile().then(p => Object.assign(store, p));
});

export const update = user => update_profile(user).then(p => Object.assign(store, p));