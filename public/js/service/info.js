
import { get_info } from "../api/info.js";
import store from "../store/info.js";

export const fetch_info = () => get_info().then(i => Object.assign(store, i));