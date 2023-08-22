import { protected_fetch } from "./util.js";

export const get_all = () => protected_fetch(`api/transaction`);

export const create = data => protected_fetch(`api/transaction/create`, {
    metod: "POST",
    body: JSON.stringify(data)
});

export const callback = data => protected_fetch(`api/transaction/callback`, {
    metod: "POST",
    body: JSON.stringify(data)
});
