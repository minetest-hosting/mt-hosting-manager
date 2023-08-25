import { get_all } from "../api/nodetype.js";

const store = Vue.reactive({
    nodetypes: []
});

export const fetch_nodetypes = () => get_all().then(nt => store.nodetypes = nt);

export const get_nodetype = id => store.nodetypes.find(nt => nt.id == id);
export const get_nodetypes = () => store.nodetypes;