import { get_all } from "../api/nodetype.js";
import store from "../store/nodetype.js";

export const fetch_nodetypes = () => get_all().then(nt => store.nodetypes = nt);

export const get_nodetype = id => store.nodetypes.find(nt => nt.id == id);