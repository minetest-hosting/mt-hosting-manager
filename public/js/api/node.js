import { protected_fetch } from "./utiljs";

export const get_all = () => protected_fetch(`api/node`);

export const create = n => protected_fetch(`api/node`, {
    metod: "POST",
    body: JSON.stringify(n)
});

export const update = n => protected_fetch(`api/node/${n.id}`, {
    metod: "POST",
    body: JSON.stringify(n)
});

export const remove = n => protected_fetch(`api/node/${n.id}`, {
    metod: "DELETE"
});