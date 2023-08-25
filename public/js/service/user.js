
import { get_profile, update_profile } from "../api/user.js";
import store from '../store/user.js';
import events, { EVENT_LOGGED_IN } from '../events.js';

events.on(EVENT_LOGGED_IN, function() {
    fetch_profile();
});

export const update = user => update_profile(user).then(p => Object.assign(store, p));

export const fetch_profile = () => get_profile().then(p => Object.assign(store, p));