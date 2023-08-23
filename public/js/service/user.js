
import { get_profile } from "../api/user.js";
import store from '../store/user.js';
import events, { EVENT_LOGGED_IN } from '../events.js';

events.on(EVENT_LOGGED_IN, function() {
    get_profile().then(p => Object.assign(store, p));
});