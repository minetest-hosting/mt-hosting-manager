import { protected_fetch } from "./protected_fetch.js";

export const get_by_id = id => protected_fetch(`api/backup/${id}`);

export const get_all = () => protected_fetch(`api/backup`);
 
export const create = b => protected_fetch(`api/backup`, {
    method: "POST",
    body: JSON.stringify(b)
});

export const remove = b => protected_fetch(`api/backup/${b.id}`, {
    method: "DELETE"
});