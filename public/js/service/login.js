
import { get_claims as fetch_claims, logout as api_logout } from '../api/login.js';
import { login as api_login } from '../api/login.js';
import events, { EVENT_LOGGED_IN } from '../events.js';

const store = Vue.reactive({
    claims: null
});

export const login = creds => api_login(creds).then(resp => {
    if (resp.status && resp.status >= 400) {
        return false;
    } else {
        return check_login()
        .then(() => true);
    }
});

export const check_login = () => fetch_claims().then(c => {
    store.claims = c;
    if (c) {
        events.emit(EVENT_LOGGED_IN, c);
    }
    return c;
});

export const logout = () => api_logout()
    .then(() => store.claims = null);

export const get_claims = () => store.claims;
export const is_logged_in = () => store.claims;
export const has_role = role => store.claims && store.claims.role == role;
