import { get_exchange_rates } from "../api/exchange_rate.js";

const store = Vue.reactive({
    rates: []
});

export const load_exchange_rates = () => get_exchange_rates().then(r => store.rates = r);

export const get_rates = () => store.rates;

export const get_rate = currency => store.rates.find(r => r.currency == currency);