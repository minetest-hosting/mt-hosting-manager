import { protected_fetch } from "./util.js";

export const get_all = () => protected_fetch(`api/nodetype`);

export const get_by_id = id => protected_fetch(`api/nodetype/${id}`);

export const add = nt => protected_fetch(`api/nodetype`, {
    method: "POST",
    body: JSON.stringify(nt)
});

export const update = nt => protected_fetch(`api/nodetype/${nt.id}`, {
    method: "POST",
    body: JSON.stringify(nt)
});

export const remove = nt => protected_fetch(`api/nodetype/${nt.id}`, {
    method: "DELETE"
});