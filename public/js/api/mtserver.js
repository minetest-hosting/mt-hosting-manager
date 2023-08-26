import { protected_fetch } from "./protected_fetch.js";

export const get_all = () => protected_fetch(`api/mtserver`);

export const get_by_id = id => protected_fetch(`api/mtserver/${id}`);

export const create = s => protected_fetch(`api/mtserver`, {
    method: "POST",
    body: JSON.stringify(s)
});

export const update = s => protected_fetch(`api/mtserver/${s.id}`, {
    method: "POST",
    body: JSON.stringify(s)
});

export const remove = s => protected_fetch(`api/mtserver/${s.id}`, {
    method: "DELETE"
});