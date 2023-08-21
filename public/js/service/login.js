
import login_store from '../store/login.js';
import { get_claims, logout as api_logout } from '../api/login.js';
import events, { EVENT_LOGGED_IN } from '../events.js';

export const check_login = () => get_claims().then(c => {
    login_store.claims = c;
    if (c) {
        events.emit(EVENT_LOGGED_IN, c);
    }
    return c;
});

export const is_logged_in = () => login_store.claims;

export const has_role = role => login_store.claims && login_store.claims.role == role;

export const logout = () => {
    return api_logout()
    .then(() => {
        login_store.claims = null;
    });
};