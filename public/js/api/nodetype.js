import { protected_fetch } from "./utiljs";

export const get_all = () => protected_fetch(`api/nodetype`);

export const add = nt => protected_fetch(`api/nodetype`, {
    metod: "POST",
    body: JSON.stringify(nt)
});

export const update = nt => protected_fetch(`api/nodetype/${nt.id}`, {
    metod: "POST",
    body: JSON.stringify(nt)
});

export const remove = nt => protected_fetch(`api/nodetype/${nt.id}`, {
    metod: "DELETE"
});