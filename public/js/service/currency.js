import { get_currency_info } from "../api/currency.js";

const store = Vue.reactive({});

export const load_currencies = () => get_currency_info().then(c => {
    Object.assign(store, c);
    const list = [];

    store.currencies
    .filter(c => store.rates.rates[c.id])
    .forEach(c => list.push({
        id: c.id,
        name: `${c.name} (${c.id})`,
        decimals: c.min_size.length-2
    }));

    store.crypto_currencies
    .filter(c => store.rates.rates[c.code])
    .forEach(c => list.push({
        id: c.code,
        name: `${c.name} (${c.code})`,
        decimals: c.exponent
    }));

    list.sort((a,b) => a.name > b.name);

    store.list = list;
});

export const get_currency_list = () => store.list;

export const get_rates = () => store.rates.rates;