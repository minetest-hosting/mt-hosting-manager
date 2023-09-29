import { get_currency_info } from "../api/currency.js";

const store = Vue.reactive({});

export const load_currencies = () => get_currency_info().then(c => Object.assign(store, c));