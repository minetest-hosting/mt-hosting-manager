
import { get_info } from "../api/info.js";

const store = Vue.reactive({});

export const fetch_info = () => get_info().then(i => Object.assign(store, i));

export const get_baseurl = () => store.baseurl;
export const get_stage = () => store.stage;
export const get_github_client_id = () => store.github_client_id;
export const get_discord_client_id = () => store.discord_client_id;
export const get_mesehub_client_id = () => store.mesehub_client_id;
export const get_hostingdomain_suffix = () => store.hostingdomain_suffix;
export const get_max_balance = () => store.max_balance;
export const get_wallee_enabled = () => store.wallee_enabled;
export const get_coinbase_enabled = () => store.coinbase_enabled;