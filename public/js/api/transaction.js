import { protected_fetch } from "./protected_fetch.js";

export const get_all = () => protected_fetch(`api/transaction`);

export const get_by_id = id => protected_fetch(`api/transaction/${id}`);

export const create = data => protected_fetch(`api/transaction/create`, {
    method: "POST",
    body: JSON.stringify(data)
});

export const check = tx => protected_fetch(`api/transaction/${tx.id}/check`);

export const search_transaction = s => protected_fetch(`api/transaction/search`, {
    method: "POST",
    body: JSON.stringify(s)
});
