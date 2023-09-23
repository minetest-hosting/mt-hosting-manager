
import { get_profile, update_profile } from "../api/profile.js";
import events, { EVENT_LOGGED_IN } from '../events.js';

const store = Vue.reactive({});

events.on(EVENT_LOGGED_IN, function() {
    fetch_profile();
});

export const update = user => update_profile(user).then(p => Object.assign(store, p));

export const fetch_profile = () => get_profile().then(p => Object.assign(store, p));

export const get_user_profile = () => store;