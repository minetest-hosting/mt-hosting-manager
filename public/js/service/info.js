
import { get_info } from "../api/info.js";

const store = Vue.reactive({});

export const fetch_info = () => get_info().then(i => Object.assign(store, i));

export const get_stage = () => store.stage;
export const get_github_client_id = () => store.github_client_id;
export const get_hostingdomain_suffix = () => store.hostingdomain_suffix;