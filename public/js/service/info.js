
import { get_info } from "../api/info.js";

const store = Vue.reactive({});

export const fetch_info = () => get_info().then(i => Object.assign(store, i));

export const get_baseurl = () => store.baseurl;
export const get_stage = () => store.stage;

export const get_github_login = () => store.github_login;
export const get_discord_login = () => store.discord_login;
export const get_mesehub_login = () => store.mesehub_login;
export const get_cdb_login = () => store.cdb_login;

export const get_hostingdomain_suffix = () => store.hostingdomain_suffix;
export const get_max_balance = () => store.max_balance;
export const get_zahlsch_enabled = () => store.zahlsch_enabled;
export const get_coinbase_enabled = () => store.coinbase_enabled;