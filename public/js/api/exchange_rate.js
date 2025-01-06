import { protected_fetch } from "./protected_fetch.js";

export const get_exchange_rates = () => protected_fetch(`api/exchange_rate`);

export const get_exchange_rate = currency => protected_fetch(`api/exchange_rate/${currency}`);

export const create_exchange_rate = er => protected_fetch(`api/exchange_rate`, {
    method: "POST",
    body: JSON.stringify(er)
});

export const update_exchange_rate = er => protected_fetch(`api/exchange_rate/${er.currency}`, {
    method: "PUT",
    body: JSON.stringify(er)
});

export const delete_exchange_rate = currency => protected_fetch(`api/exchange_rate/${currency}`, {
    method: "DELETE"
});
